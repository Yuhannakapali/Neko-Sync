# Database Schema

This document describes the PostgreSQL database schema for the Neko-Sync backend.

## Overview

The database follows a normalized design with clear relationships between entities. All tables use UUIDs as primary keys for better scalability and security.

## Tables

### Users and Authentication

#### users
Core user accounts and authentication data.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    avatar_url TEXT,
    role VARCHAR(20) DEFAULT 'user' CHECK (role IN ('user', 'admin', 'moderator')),
    is_verified BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_created_at ON users(created_at);
```

#### user_profiles
Extended user profile information.

```sql
CREATE TABLE user_profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    about TEXT,
    location VARCHAR(255),
    website VARCHAR(255),
    banner_url TEXT,
    birthdate DATE,
    privacy_settings JSONB DEFAULT '{}',
    preferences JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### user_devices
User's connected devices for content synchronization.

```sql
CREATE TABLE user_devices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_name VARCHAR(255) NOT NULL,
    platform VARCHAR(50) NOT NULL CHECK (platform IN ('web', 'desktop', 'mobile', 'tablet')),
    device_token VARCHAR(255),
    websocket_id VARCHAR(255),
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_user_devices_user_id ON user_devices(user_id);
CREATE INDEX idx_user_devices_active ON user_devices(is_active);
```

### Content Management

#### content_series
Content series/franchises (anime series, manga series, etc.).

```sql
CREATE TABLE content_series (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(500) NOT NULL,
    description TEXT,
    cover_url TEXT,
    type VARCHAR(20) NOT NULL CHECK (type IN ('anime', 'manga', 'movie', 'music')),
    status VARCHAR(20) DEFAULT 'ongoing' CHECK (status IN ('ongoing', 'completed', 'upcoming', 'cancelled')),
    release_date DATE,
    end_date DATE,
    rating DECIMAL(3,2) DEFAULT 0.0,
    total_episodes INTEGER,
    total_chapters INTEGER,
    external_ids JSONB DEFAULT '{}', -- MAL, TMDB, etc.
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_content_series_type ON content_series(type);
CREATE INDEX idx_content_series_status ON content_series(status);
CREATE INDEX idx_content_series_rating ON content_series(rating);
CREATE INDEX idx_content_series_release_date ON content_series(release_date);
```

#### content
Individual content items (episodes, chapters, movies, songs).

```sql
CREATE TABLE content (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    series_id UUID REFERENCES content_series(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    description TEXT,
    type VARCHAR(20) NOT NULL CHECK (type IN ('episode', 'chapter', 'movie', 'song')),
    episode_number INTEGER,
    chapter_number INTEGER,
    season_number INTEGER,
    duration_minutes INTEGER,
    pages INTEGER,
    release_date DATE,
    thumbnail_url TEXT,
    file_url TEXT,
    external_ids JSONB DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    view_count BIGINT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_content_series_id ON content(series_id);
CREATE INDEX idx_content_type ON content(type);
CREATE INDEX idx_content_episode_number ON content(episode_number);
CREATE INDEX idx_content_release_date ON content(release_date);
```

#### genres
Content genres and categories.

```sql
CREATE TABLE genres (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    color VARCHAR(7), -- Hex color code
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### content_genres
Many-to-many relationship between content and genres.

```sql
CREATE TABLE content_genres (
    content_id UUID REFERENCES content_series(id) ON DELETE CASCADE,
    genre_id UUID REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (content_id, genre_id)
);
```

#### tags
Flexible tagging system for content.

```sql
CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    category VARCHAR(50), -- 'mood', 'theme', 'setting', etc.
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### content_tags
Many-to-many relationship between content and tags.

```sql
CREATE TABLE content_tags (
    content_id UUID REFERENCES content_series(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (content_id, tag_id)
);
```

### Watch Parties

#### watch_parties
Synchronized viewing sessions.

```sql
CREATE TABLE watch_parties (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    host_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content_series(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'waiting' CHECK (status IN ('waiting', 'active', 'paused', 'ended')),
    max_participants INTEGER DEFAULT 10,
    current_participants INTEGER DEFAULT 1,
    is_public BOOLEAN DEFAULT TRUE,
    join_code VARCHAR(10) UNIQUE,
    current_episode INTEGER DEFAULT 1,
    current_timestamp BIGINT DEFAULT 0, -- milliseconds
    started_at TIMESTAMP WITH TIME ZONE,
    ended_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_watch_parties_host_id ON watch_parties(host_id);
CREATE INDEX idx_watch_parties_content_id ON watch_parties(content_id);
CREATE INDEX idx_watch_parties_status ON watch_parties(status);
CREATE INDEX idx_watch_parties_join_code ON watch_parties(join_code);
CREATE INDEX idx_watch_parties_is_public ON watch_parties(is_public);
```

#### watch_party_participants
Users participating in watch parties.

```sql
CREATE TABLE watch_party_participants (
    party_id UUID REFERENCES watch_parties(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    is_active BOOLEAN DEFAULT TRUE,
    PRIMARY KEY (party_id, user_id)
);

CREATE INDEX idx_watch_party_participants_user_id ON watch_party_participants(user_id);
```

#### party_messages
Chat messages within watch parties.

```sql
CREATE TABLE party_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    party_id UUID NOT NULL REFERENCES watch_parties(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message TEXT NOT NULL,
    timestamp_ms BIGINT, -- Video timestamp when message was sent
    message_type VARCHAR(20) DEFAULT 'chat' CHECK (message_type IN ('chat', 'system', 'reaction')),
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_party_messages_party_id ON party_messages(party_id);
CREATE INDEX idx_party_messages_created_at ON party_messages(created_at);
```

### User History and Progress

#### watch_history
User's watch progress and history.

```sql
CREATE TABLE watch_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    series_id UUID NOT NULL REFERENCES content_series(id) ON DELETE CASCADE,
    episode_number INTEGER,
    timestamp_ms BIGINT DEFAULT 0,
    duration_ms BIGINT,
    completed BOOLEAN DEFAULT FALSE,
    watch_count INTEGER DEFAULT 1,
    first_watched_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_watched_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_watch_history_user_content ON watch_history(user_id, content_id);
CREATE INDEX idx_watch_history_user_id ON watch_history(user_id);
CREATE INDEX idx_watch_history_series_id ON watch_history(series_id);
CREATE INDEX idx_watch_history_last_watched ON watch_history(last_watched_at);
```

#### read_history
User's reading progress for manga/novels.

```sql
CREATE TABLE read_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    series_id UUID NOT NULL REFERENCES content_series(id) ON DELETE CASCADE,
    chapter_number INTEGER,
    page_number INTEGER DEFAULT 0,
    completed BOOLEAN DEFAULT FALSE,
    read_count INTEGER DEFAULT 1,
    first_read_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_read_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_read_history_user_content ON read_history(user_id, content_id);
CREATE INDEX idx_read_history_user_id ON read_history(user_id);
CREATE INDEX idx_read_history_series_id ON read_history(series_id);
```

#### favorites
User's favorite content.

```sql
CREATE TABLE favorites (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    content_id UUID REFERENCES content_series(id) ON DELETE CASCADE,
    favorited_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, content_id)
);

CREATE INDEX idx_favorites_user_id ON favorites(user_id);
CREATE INDEX idx_favorites_favorited_at ON favorites(favorited_at);
```

### Social Features

#### user_follows
User follow relationships.

```sql
CREATE TABLE user_follows (
    follower_id UUID REFERENCES users(id) ON DELETE CASCADE,
    following_id UUID REFERENCES users(id) ON DELETE CASCADE,
    followed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (follower_id, following_id),
    CHECK (follower_id != following_id)
);

CREATE INDEX idx_user_follows_follower_id ON user_follows(follower_id);
CREATE INDEX idx_user_follows_following_id ON user_follows(following_id);
```

#### discussions
Content-related discussions and forums.

```sql
CREATE TABLE discussions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content_id UUID REFERENCES content_series(id) ON DELETE SET NULL,
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    view_count BIGINT DEFAULT 0,
    reply_count INTEGER DEFAULT 0,
    is_pinned BOOLEAN DEFAULT FALSE,
    is_locked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_discussions_user_id ON discussions(user_id);
CREATE INDEX idx_discussions_content_id ON discussions(content_id);
CREATE INDEX idx_discussions_created_at ON discussions(created_at);
```

#### discussion_posts
Replies to discussions.

```sql
CREATE TABLE discussion_posts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    discussion_id UUID NOT NULL REFERENCES discussions(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    parent_id UUID REFERENCES discussion_posts(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_discussion_posts_discussion_id ON discussion_posts(discussion_id);
CREATE INDEX idx_discussion_posts_user_id ON discussion_posts(user_id);
CREATE INDEX idx_discussion_posts_parent_id ON discussion_posts(parent_id);
```

#### content_reviews
User reviews and ratings for content.

```sql
CREATE TABLE content_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content_series(id) ON DELETE CASCADE,
    rating INTEGER CHECK (rating >= 1 AND rating <= 10),
    title VARCHAR(255),
    content TEXT,
    is_spoiler BOOLEAN DEFAULT FALSE,
    helpful_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, content_id)
);

CREATE INDEX idx_content_reviews_content_id ON content_reviews(content_id);
CREATE INDEX idx_content_reviews_rating ON content_reviews(rating);
```

#### content_comments
Comments on content (episodes, chapters, etc.).

```sql
CREATE TABLE content_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    comment TEXT NOT NULL,
    parent_id UUID REFERENCES content_comments(id) ON DELETE CASCADE,
    timestamp_ms BIGINT, -- Video/reading timestamp for the comment
    is_spoiler BOOLEAN DEFAULT FALSE,
    like_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_content_comments_content_id ON content_comments(content_id);
CREATE INDEX idx_content_comments_user_id ON content_comments(user_id);
CREATE INDEX idx_content_comments_parent_id ON content_comments(parent_id);
```

### Notifications

#### notifications
User notifications system.

```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL CHECK (type IN ('follow', 'like', 'comment', 'party_invite', 'content_update')),
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSONB DEFAULT '{}', -- Additional notification data
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
```

### Device Transfer

#### device_transfers
Content transfer between user devices.

```sql
CREATE TABLE device_transfers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    from_device_id UUID NOT NULL REFERENCES user_devices(id) ON DELETE CASCADE,
    to_device_id UUID NOT NULL REFERENCES user_devices(id) ON DELETE CASCADE,
    content_id UUID NOT NULL REFERENCES content(id) ON DELETE CASCADE,
    timestamp_ms BIGINT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'in_progress', 'completed', 'failed')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_device_transfers_user_id ON device_transfers(user_id);
CREATE INDEX idx_device_transfers_status ON device_transfers(status);
```

## Views

### Popular Content View
```sql
CREATE VIEW popular_content AS
SELECT 
    cs.*,
    COALESCE(wh.watch_count, 0) as total_watches,
    COALESCE(f.favorite_count, 0) as favorite_count,
    COALESCE(cr.avg_rating, 0) as average_rating,
    COALESCE(cr.review_count, 0) as review_count
FROM content_series cs
LEFT JOIN (
    SELECT series_id, SUM(watch_count) as watch_count
    FROM watch_history
    GROUP BY series_id
) wh ON cs.id = wh.series_id
LEFT JOIN (
    SELECT content_id, COUNT(*) as favorite_count
    FROM favorites
    GROUP BY content_id
) f ON cs.id = f.content_id
LEFT JOIN (
    SELECT content_id, AVG(rating) as avg_rating, COUNT(*) as review_count
    FROM content_reviews
    GROUP BY content_id
) cr ON cs.id = cr.content_id;
```

## Functions and Triggers

### Update timestamps
```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply to all tables with updated_at column
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE ON user_profiles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_content_series_updated_at BEFORE UPDATE ON content_series FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_content_updated_at BEFORE UPDATE ON content FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_watch_parties_updated_at BEFORE UPDATE ON watch_parties FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

### Generate join codes for watch parties
```sql
CREATE OR REPLACE FUNCTION generate_join_code()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.join_code IS NULL THEN
        NEW.join_code = upper(substr(md5(random()::text), 1, 6));
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER generate_watch_party_join_code 
    BEFORE INSERT ON watch_parties 
    FOR EACH ROW EXECUTE FUNCTION generate_join_code();
```

## Indexes for Performance

```sql
-- Composite indexes for common queries
CREATE INDEX idx_watch_history_user_series ON watch_history(user_id, series_id, last_watched_at);
CREATE INDEX idx_content_type_status ON content_series(type, status, rating DESC);
CREATE INDEX idx_watch_parties_public_active ON watch_parties(is_public, status) WHERE is_public = true AND status IN ('waiting', 'active');
CREATE INDEX idx_notifications_user_unread ON notifications(user_id, created_at DESC) WHERE is_read = false;

-- Full-text search indexes
CREATE INDEX idx_content_series_title_search ON content_series USING gin(to_tsvector('english', title));
CREATE INDEX idx_content_series_description_search ON content_series USING gin(to_tsvector('english', description));
```

## Data Migration

For initial setup, run the migrations in this order:

1. Create tables (in dependency order)
2. Create indexes
3. Create views
4. Create functions and triggers
5. Insert seed data (genres, default admin user, etc.)

## Backup Strategy

- **Daily backups** of the entire database
- **Point-in-time recovery** enabled
- **Separate backups** for user data and content metadata
- **Test restore procedures** monthly

## Performance Considerations

- Use **connection pooling** (pgbouncer recommended)
- Regular **VACUUM** and **ANALYZE** operations
- Monitor **slow queries** and optimize as needed
- Consider **read replicas** for heavy read workloads
- **Partition** large tables like `watch_history` by date if needed
