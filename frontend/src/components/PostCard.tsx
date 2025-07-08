import React from 'react';
import { Post } from '../types';

interface PostCardProps {
  post: Post;
  onLike: (postId: string) => void;
  onRepost: (postId: string) => void;
  onDelete?: (postId: string) => void;
  currentUserId?: string;
}

const PostCard: React.FC<PostCardProps> = ({ 
  post, 
  onLike, 
  onRepost, 
  onDelete,
  currentUserId 
}) => {
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ja-JP', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <div style={cardStyle}>
      <div style={headerStyle}>
        <div style={avatarStyle}>
          {post.author.avatarUrl ? (
            <img src={post.author.avatarUrl} alt="avatar" style={avatarImageStyle} />
          ) : (
            <div style={avatarPlaceholderStyle}>
              {post.author.displayName.charAt(0)}
            </div>
          )}
        </div>
        
        <div style={authorInfoStyle}>
          <span style={displayNameStyle}>{post.author.displayName}</span>
          <span style={usernameStyle}>@{post.author.username}</span>
          <span style={dateStyle}>・{formatDate(post.createdAt)}</span>
        </div>
        
        {currentUserId === post.author.id && onDelete && (
          <button 
            onClick={() => onDelete(post.id)}
            style={deleteButtonStyle}
          >
            削除
          </button>
        )}
      </div>
      
      <div style={contentStyle}>
        <p>{post.content}</p>
        
        {post.imageUrls && post.imageUrls.length > 0 && (
          <div style={imagesStyle}>
            {post.imageUrls.map((url, index) => (
              <img 
                key={index} 
                src={url} 
                alt={`Post image ${index + 1}`}
                style={imageStyle}
              />
            ))}
          </div>
        )}
      </div>
      
      <div style={actionsStyle}>
        <button 
          onClick={() => onLike(post.id)}
          style={{
            ...actionButtonStyle,
            color: post.isLiked ? '#e0245e' : '#657786',
          }}
        >
          ♥ {post.likesCount}
        </button>
        
        <button 
          onClick={() => onRepost(post.id)}
          style={{
            ...actionButtonStyle,
            color: post.isReposted ? '#17bf63' : '#657786',
          }}
        >
          ↻ {post.repostsCount}
        </button>
        
        <button style={actionButtonStyle}>
          💬 {post.commentsCount}
        </button>
      </div>
    </div>
  );
};

const cardStyle: React.CSSProperties = {
  backgroundColor: 'white',
  border: '1px solid #e1e8ed',
  borderRadius: '8px',
  padding: '1rem',
  marginBottom: '1rem',
};

const headerStyle: React.CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  marginBottom: '0.5rem',
};

const avatarStyle: React.CSSProperties = {
  marginRight: '0.75rem',
};

const avatarImageStyle: React.CSSProperties = {
  width: '40px',
  height: '40px',
  borderRadius: '50%',
  objectFit: 'cover',
};

const avatarPlaceholderStyle: React.CSSProperties = {
  width: '40px',
  height: '40px',
  borderRadius: '50%',
  backgroundColor: '#1da1f2',
  color: 'white',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  fontWeight: 'bold',
};

const authorInfoStyle: React.CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  gap: '0.5rem',
  flex: 1,
};

const displayNameStyle: React.CSSProperties = {
  fontWeight: 'bold',
  color: '#14171a',
};

const usernameStyle: React.CSSProperties = {
  color: '#657786',
};

const dateStyle: React.CSSProperties = {
  color: '#657786',
  fontSize: '0.9rem',
};

const deleteButtonStyle: React.CSSProperties = {
  backgroundColor: 'transparent',
  border: 'none',
  color: '#e0245e',
  cursor: 'pointer',
  fontSize: '0.9rem',
};

const contentStyle: React.CSSProperties = {
  marginBottom: '1rem',
};

const imagesStyle: React.CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(200px, 1fr))',
  gap: '0.5rem',
  marginTop: '0.5rem',
};

const imageStyle: React.CSSProperties = {
  width: '100%',
  height: '200px',
  objectFit: 'cover',
  borderRadius: '8px',
};

const actionsStyle: React.CSSProperties = {
  display: 'flex',
  gap: '2rem',
  paddingTop: '0.5rem',
  borderTop: '1px solid #e1e8ed',
};

const actionButtonStyle: React.CSSProperties = {
  backgroundColor: 'transparent',
  border: 'none',
  color: '#657786',
  cursor: 'pointer',
  display: 'flex',
  alignItems: 'center',
  gap: '0.5rem',
  fontSize: '0.9rem',
};

export default PostCard;