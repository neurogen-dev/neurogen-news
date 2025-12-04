-- Migration: Extended Schema
-- Adds new fields and tables for full feature support

-- ============================================
-- Users table extensions
-- ============================================
ALTER TABLE users ADD COLUMN IF NOT EXISTS cover_url VARCHAR(500);
ALTER TABLE users ADD COLUMN IF NOT EXISTS location VARCHAR(100);
ALTER TABLE users ADD COLUMN IF NOT EXISTS website VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS telegram VARCHAR(50);
ALTER TABLE users ADD COLUMN IF NOT EXISTS github VARCHAR(50);

-- ============================================
-- Categories table extensions
-- ============================================
ALTER TABLE categories ADD COLUMN IF NOT EXISTS color VARCHAR(7);
ALTER TABLE categories ADD COLUMN IF NOT EXISTS is_official BOOLEAN DEFAULT false;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS parent_id UUID REFERENCES categories(id);
ALTER TABLE categories ADD COLUMN IF NOT EXISTS subscriber_count INT DEFAULT 0;
ALTER TABLE categories ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT NOW();

-- ============================================
-- Achievements table
-- ============================================
CREATE TABLE IF NOT EXISTS achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    icon VARCHAR(50) NOT NULL,
    points INT DEFAULT 0,
    category VARCHAR(50) DEFAULT 'general', -- author, commentator, social, special
    is_hidden BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- User achievements
-- ============================================
CREATE TABLE IF NOT EXISTS user_achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id UUID NOT NULL REFERENCES achievements(id) ON DELETE CASCADE,
    awarded_at TIMESTAMP DEFAULT NOW(),
    UNIQUE(user_id, achievement_id)
);

-- ============================================
-- Bookmark folders
-- ============================================
CREATE TABLE IF NOT EXISTS bookmark_folders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Add folder_id to bookmarks
ALTER TABLE bookmarks ADD COLUMN IF NOT EXISTS folder_id UUID REFERENCES bookmark_folders(id) ON DELETE SET NULL;

-- ============================================
-- Drafts table extensions
-- ============================================
ALTER TABLE drafts ADD COLUMN IF NOT EXISTS article_id UUID REFERENCES articles(id) ON DELETE CASCADE;
ALTER TABLE drafts ADD COLUMN IF NOT EXISTS is_auto_save BOOLEAN DEFAULT false;
ALTER TABLE drafts ADD COLUMN IF NOT EXISTS cover_image_url VARCHAR(500);

-- ============================================
-- Notifications table extensions
-- ============================================
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS article_id UUID REFERENCES articles(id) ON DELETE CASCADE;
ALTER TABLE notifications ADD COLUMN IF NOT EXISTS comment_id UUID REFERENCES comments(id) ON DELETE CASCADE;

-- ============================================
-- Moderation tables
-- ============================================
CREATE TABLE IF NOT EXISTS moderation_actions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    moderator_id UUID NOT NULL REFERENCES users(id),
    target_type VARCHAR(20) NOT NULL, -- article, comment, user
    target_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL, -- approve, reject, delete, ban, warn
    reason TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_bans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    banned_by UUID NOT NULL REFERENCES users(id),
    reason TEXT NOT NULL,
    expires_at TIMESTAMP,
    is_permanent BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Article versions (for edit history)
-- ============================================
CREATE TABLE IF NOT EXISTS article_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    article_id UUID NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    edited_by UUID NOT NULL REFERENCES users(id),
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW()
);

-- ============================================
-- Indexes
-- ============================================
CREATE INDEX IF NOT EXISTS idx_user_achievements_user_id ON user_achievements(user_id);
CREATE INDEX IF NOT EXISTS idx_bookmark_folders_user_id ON bookmark_folders(user_id);
CREATE INDEX IF NOT EXISTS idx_bookmarks_folder_id ON bookmarks(folder_id);
CREATE INDEX IF NOT EXISTS idx_moderation_actions_target ON moderation_actions(target_type, target_id);
CREATE INDEX IF NOT EXISTS idx_user_bans_user_id ON user_bans(user_id);
CREATE INDEX IF NOT EXISTS idx_article_versions_article_id ON article_versions(article_id);
CREATE INDEX IF NOT EXISTS idx_notifications_article_id ON notifications(article_id);
CREATE INDEX IF NOT EXISTS idx_notifications_comment_id ON notifications(comment_id);
CREATE INDEX IF NOT EXISTS idx_drafts_auto_save ON drafts(user_id, article_id) WHERE is_auto_save = true;

-- ============================================
-- Seed default achievements
-- ============================================
INSERT INTO achievements (name, slug, description, icon, points, category) VALUES
    ('Первая статья', 'first-article', 'Опубликуйте свою первую статью', '📝', 10, 'author'),
    ('Плодовитый автор', 'prolific-writer', 'Опубликуйте 10 статей', '✍️', 50, 'author'),
    ('Сотня статей', 'author-100', 'Опубликуйте 100 статей', '📚', 200, 'author'),
    ('Первый комментарий', 'first-comment', 'Оставьте свой первый комментарий', '💬', 5, 'commentator'),
    ('Активный комментатор', 'active-commenter', 'Оставьте 50 комментариев', '🗨️', 30, 'commentator'),
    ('Сотня комментариев', 'commentator-100', 'Оставьте 100 комментариев', '💭', 100, 'commentator'),
    ('Популярный автор', 'popular-author', 'Получите 1000 просмотров на статьи', '👀', 25, 'social'),
    ('Вирусный автор', 'viral-author', 'Получите 10000 просмотров на статьи', '🔥', 100, 'social'),
    ('Первый подписчик', 'first-follower', 'Получите первого подписчика', '👤', 10, 'social'),
    ('Влиятельный', 'influencer', 'Получите 100 подписчиков', '⭐', 75, 'social'),
    ('Знаменитость', 'celebrity', 'Получите 1000 подписчиков', '🌟', 250, 'social')
ON CONFLICT (slug) DO NOTHING;

-- ============================================
-- Update triggers
-- ============================================
CREATE OR REPLACE FUNCTION update_category_subscriber_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE categories SET subscriber_count = subscriber_count + 1 WHERE id = NEW.category_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE categories SET subscriber_count = subscriber_count - 1 WHERE id = OLD.category_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_category_subscriber_count ON category_subscriptions;
CREATE TRIGGER trigger_update_category_subscriber_count
AFTER INSERT OR DELETE ON category_subscriptions
FOR EACH ROW EXECUTE FUNCTION update_category_subscriber_count();

