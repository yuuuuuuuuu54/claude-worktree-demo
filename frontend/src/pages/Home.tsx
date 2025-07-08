import React, { useState, useEffect } from 'react';
import { Post, User, CreatePostData } from '../types';
import { apiClient } from '../api/client';
import PostForm from '../components/PostForm';
import PostCard from '../components/PostCard';

interface HomeProps {
  user: User | null;
}

const Home: React.FC<HomeProps> = ({ user }) => {
  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [posting, setPosting] = useState(false);

  useEffect(() => {
    fetchPosts();
  }, []);

  const fetchPosts = async () => {
    try {
      const response = await apiClient.getPosts();
      setPosts(response.data);
    } catch (error) {
      console.error('Failed to fetch posts:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleCreatePost = async (data: CreatePostData) => {
    if (!user) return;
    
    setPosting(true);
    try {
      const response = await apiClient.createPost(data);
      setPosts(prev => [response.data, ...prev]);
    } catch (error) {
      console.error('Failed to create post:', error);
    } finally {
      setPosting(false);
    }
  };

  const handleLike = async (postId: string) => {
    if (!user) return;
    
    try {
      const post = posts.find(p => p.id === postId);
      if (!post) return;

      if (post.isLiked) {
        await apiClient.unlikePost(postId);
      } else {
        await apiClient.likePost(postId);
      }

      setPosts(prev => prev.map(p => 
        p.id === postId 
          ? { 
              ...p, 
              isLiked: !p.isLiked, 
              likesCount: p.isLiked ? p.likesCount - 1 : p.likesCount + 1 
            }
          : p
      ));
    } catch (error) {
      console.error('Failed to like/unlike post:', error);
    }
  };

  const handleRepost = async (postId: string) => {
    if (!user) return;
    
    try {
      const post = posts.find(p => p.id === postId);
      if (!post) return;

      if (post.isReposted) {
        await apiClient.unrepost(postId);
      } else {
        await apiClient.repost(postId);
      }

      setPosts(prev => prev.map(p => 
        p.id === postId 
          ? { 
              ...p, 
              isReposted: !p.isReposted, 
              repostsCount: p.isReposted ? p.repostsCount - 1 : p.repostsCount + 1 
            }
          : p
      ));
    } catch (error) {
      console.error('Failed to repost/unrepost:', error);
    }
  };

  const handleDelete = async (postId: string) => {
    if (!user) return;
    
    try {
      await apiClient.deletePost(postId);
      setPosts(prev => prev.filter(p => p.id !== postId));
    } catch (error) {
      console.error('Failed to delete post:', error);
    }
  };

  if (loading) {
    return (
      <div style={loadingStyle}>
        <p>投稿を読み込んでいます...</p>
      </div>
    );
  }

  return (
    <div style={containerStyle}>
      <div style={contentStyle}>
        <h2 style={titleStyle}>ホーム</h2>
        
        {user && (
          <PostForm 
            onSubmit={handleCreatePost}
            loading={posting}
          />
        )}
        
        <div style={timelineStyle}>
          {posts.length === 0 ? (
            <div style={emptyStateStyle}>
              <p>まだ投稿がありません</p>
              {user && <p>最初の投稿をしてみましょう！</p>}
            </div>
          ) : (
            posts.map(post => (
              <PostCard
                key={post.id}
                post={post}
                onLike={handleLike}
                onRepost={handleRepost}
                onDelete={handleDelete}
                currentUserId={user?.id}
              />
            ))
          )}
        </div>
      </div>
    </div>
  );
};

const containerStyle: React.CSSProperties = {
  maxWidth: '600px',
  margin: '0 auto',
  padding: '1rem',
};

const contentStyle: React.CSSProperties = {
  backgroundColor: '#f5f8fa',
  minHeight: '100vh',
};

const titleStyle: React.CSSProperties = {
  fontSize: '1.5rem',
  fontWeight: 'bold',
  marginBottom: '1rem',
  color: '#14171a',
};

const timelineStyle: React.CSSProperties = {
  marginTop: '1rem',
};

const loadingStyle: React.CSSProperties = {
  display: 'flex',
  justifyContent: 'center',
  alignItems: 'center',
  height: '200px',
  color: '#657786',
};

const emptyStateStyle: React.CSSProperties = {
  textAlign: 'center',
  padding: '2rem',
  color: '#657786',
  backgroundColor: 'white',
  borderRadius: '8px',
  border: '1px solid #e1e8ed',
};

export default Home;