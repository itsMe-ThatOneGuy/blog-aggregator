package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/commands"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/database"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/rss"
	"github.com/itsMe-ThatOneGuy/blog-aggregator/internal/state"
)

func Agg(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("requires a time between requests\n")
	}

	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration: %w\n", err)
	}

	fmt.Printf("collecting feeds every %s...\n", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		ScrapeFeeds(s)
	}
}

func ScrapeFeeds(s *state.State) error {
	nextFeed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("Couldn't get next feed to fetch: %w\n", err)
	}

	_, err = s.DB.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return fmt.Errorf("Couldn't mark feed %s as fetched: %w\n", nextFeed.Name, err)
	}

	feed, err := rss.FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return fmt.Errorf("Couldn't collect feed %s: %w\n", nextFeed.Name, err)
	}

	for _, item := range feed.Channel.Item {
		fmt.Printf("Posts: %s\n", item.Title)
	}

	return nil
}

func Reset(s *state.State, cmd commands.Command) error {
	err := s.DB.ResetDB(context.Background())
	if err != nil {
		return fmt.Errorf("issue deleting users: %w", err)
	}
	fmt.Println("Database reset")

	return nil
}

func HandlerAddFeed(s *state.State, cmd commands.Command, user database.User) error {
	if len(cmd.Args) < 2 {
		return fmt.Errorf("add feed handler expects a name and url")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	}

	newFeed, err := s.DB.CreateFeed(context.Background(), feed)
	if err != nil {
		return fmt.Errorf("issue creating feed")
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.DB.CreateFeedFollow(context.Background(), feedFollow)
	if err != nil {
		return err
	}

	fmt.Printf(url)
	printFeed(newFeed)

	return nil
}

func HandlerListFeeds(s *state.State, cmd commands.Command) error {
	feeds, err := s.DB.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("issue retrieving feeds")
	}

	for _, feed := range feeds {
		user, err := s.DB.GetUserByID(context.Background(), feed.UserID)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", feed.Name)
		fmt.Printf(" - URL: %s\n", feed.Url)
		fmt.Printf(" - Added By: %s\n", user.Name)
	}

	return nil
}

func HandlerFollow(s *state.State, cmd commands.Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("follow requires a url")
	}
	url := cmd.Args[0]

	feed, err := s.DB.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	feedToFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	_, err = s.DB.CreateFeedFollow(context.Background(), feedToFollow)
	if err != nil {
		return fmt.Errorf("issue following feed")
	}

	fmt.Println("User")
	printUser(user)
	fmt.Println("feed")
	printFeed(feed)

	return nil
}

func HandlerUnfollow(s *state.State, cmd commands.Command, user database.User) error {
	if len(cmd.Args) < 1 {
		return fmt.Errorf("unfollow requires a feed url")
	}

	url := cmd.Args[0]

	feed, err := s.DB.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("issue getting feed: %w", err)
	}

	deleteFeed := database.DeleteFeedFollowParams{
		ID:  user.ID,
		Url: url,
	}

	err = s.DB.DeleteFeedFollow(context.Background(), deleteFeed)
	if err != nil {
		return fmt.Errorf("issue deleting follow: %w", err)
	}

	fmt.Printf("%s unfollowd %s\n", user.Name, feed.Name)

	return nil
}

func HandlerListFollowing(s *state.State, cmd commands.Command, user database.User) error {
	follows, err := s.DB.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return err
	}

	fmt.Printf("User: %v\n", user.Name)
	fmt.Println("Feeds:")
	for _, feed := range follows {
		fmt.Printf(" - %v\n", feed.FeedName)
	}

	return nil
}

func HandlerLogin(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("login handler expects a username argument")
	}
	name := cmd.Args[0]

	_, err := s.DB.GetUser(context.Background(), name)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("can't login user: %w\n", err)
		}

		return err
	}

	s.ConfigPointer.SetUser(name)
	fmt.Printf("User has been set: %s\n", name)

	return nil
}

func HandlerRegister(s *state.State, cmd commands.Command) error {
	if len(cmd.Args) == 0 {
		return errors.New("register handler expects a username")
	}
	name := cmd.Args[0]

	_, err := s.DB.GetUser(context.Background(), name)
	if err == nil {
		fmt.Printf("user '%s' already exists\n", name)
		os.Exit(1)
	} else if err != sql.ErrNoRows {
		return err
	}

	user := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name,
	}

	newuser, err := s.DB.CreateUser(context.Background(), user)
	if err != nil {
		return fmt.Errorf("Issue creating user")
	}
	s.ConfigPointer.SetUser(name)

	fmt.Println("User has been created")
	printUser(newuser)

	return nil
}

func HandlerGetUsers(s *state.State, cmd commands.Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("issue retrieving users")
	}

	for _, user := range users {
		current := s.ConfigPointer.User
		if user.Name == current {
			fmt.Printf("* %v (current)\n", user.Name)
			continue
		}
		fmt.Printf("* %v\n", user.Name)
	}

	return nil
}

func getCurrentUser(s *state.State) (database.User, error) {
	cfgUser := s.ConfigPointer.User
	user, err := s.DB.GetUser(context.Background(), cfgUser)
	if err != nil {
		if err != sql.ErrNoRows {
			return database.User{}, err
		}
		return database.User{}, fmt.Errorf("User not found")
	}

	return user, nil
}

func printUser(user database.User) {
	fmt.Printf("* ID:   %v\n", user.ID)
	fmt.Printf("* Name: %v\n", user.Name)
}

func printFeed(feed database.Feed) {
	fmt.Printf("* ID:   %v\n", feed.ID)
	fmt.Printf("* Name: %v\n", feed.Name)
	fmt.Printf("* URL: %v\n", feed.Url)
	fmt.Printf("* UserID: %v\n", feed.UserID)
}
