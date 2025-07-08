import React, { useState } from 'react';
import { CreatePostData } from '../types';

interface PostFormProps {
  onSubmit: (data: CreatePostData) => void;
  loading?: boolean;
}

const PostForm: React.FC<PostFormProps> = ({ onSubmit, loading = false }) => {
  const [content, setContent] = useState('');
  const [imageFiles, setImageFiles] = useState<File[]>([]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (content.trim()) {
      onSubmit({ content, imageFiles: imageFiles.length > 0 ? imageFiles : undefined });
      setContent('');
      setImageFiles([]);
    }
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files) {
      const files = Array.from(e.target.files).slice(0, 4);
      setImageFiles(files);
    }
  };

  const removeImage = (index: number) => {
    setImageFiles(prev => prev.filter((_, i) => i !== index));
  };

  const isValidPost = content.trim().length > 0 && content.length <= 280;

  return (
    <form onSubmit={handleSubmit} style={formStyle}>
      <div style={inputAreaStyle}>
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder="今何してる？"
          style={textareaStyle}
          rows={3}
          maxLength={280}
        />
        
        <div style={characterCountStyle}>
          <span style={{ color: content.length > 280 ? '#e0245e' : '#657786' }}>
            {content.length}/280
          </span>
        </div>
      </div>
      
      {imageFiles.length > 0 && (
        <div style={imagePreviewStyle}>
          {imageFiles.map((file, index) => (
            <div key={index} style={imageItemStyle}>
              <img 
                src={URL.createObjectURL(file)} 
                alt={`Preview ${index + 1}`}
                style={previewImageStyle}
              />
              <button 
                type="button"
                onClick={() => removeImage(index)}
                style={removeImageButtonStyle}
              >
                ×
              </button>
            </div>
          ))}
        </div>
      )}
      
      <div style={footerStyle}>
        <label style={imageButtonStyle}>
          📷 画像を追加
          <input
            type="file"
            accept="image/*"
            multiple
            onChange={handleImageChange}
            style={hiddenInputStyle}
            disabled={imageFiles.length >= 4}
          />
        </label>
        
        <button 
          type="submit" 
          disabled={!isValidPost || loading}
          style={{
            ...submitButtonStyle,
            opacity: isValidPost && !loading ? 1 : 0.6,
          }}
        >
          {loading ? '投稿中...' : '投稿'}
        </button>
      </div>
    </form>
  );
};

const formStyle: React.CSSProperties = {
  backgroundColor: 'white',
  border: '1px solid #e1e8ed',
  borderRadius: '8px',
  padding: '1rem',
  marginBottom: '1rem',
};

const inputAreaStyle: React.CSSProperties = {
  marginBottom: '1rem',
};

const textareaStyle: React.CSSProperties = {
  width: '100%',
  border: 'none',
  outline: 'none',
  resize: 'vertical',
  fontSize: '1.1rem',
  fontFamily: 'inherit',
  padding: '0.5rem',
};

const characterCountStyle: React.CSSProperties = {
  textAlign: 'right',
  fontSize: '0.9rem',
  marginTop: '0.5rem',
};

const imagePreviewStyle: React.CSSProperties = {
  display: 'grid',
  gridTemplateColumns: 'repeat(auto-fit, minmax(150px, 1fr))',
  gap: '0.5rem',
  marginBottom: '1rem',
};

const imageItemStyle: React.CSSProperties = {
  position: 'relative',
  display: 'inline-block',
};

const previewImageStyle: React.CSSProperties = {
  width: '100%',
  height: '120px',
  objectFit: 'cover',
  borderRadius: '8px',
};

const removeImageButtonStyle: React.CSSProperties = {
  position: 'absolute',
  top: '5px',
  right: '5px',
  backgroundColor: 'rgba(0, 0, 0, 0.7)',
  color: 'white',
  border: 'none',
  borderRadius: '50%',
  width: '24px',
  height: '24px',
  cursor: 'pointer',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
};

const footerStyle: React.CSSProperties = {
  display: 'flex',
  justifyContent: 'space-between',
  alignItems: 'center',
  borderTop: '1px solid #e1e8ed',
  paddingTop: '0.5rem',
};

const imageButtonStyle: React.CSSProperties = {
  backgroundColor: 'transparent',
  border: 'none',
  color: '#1da1f2',
  cursor: 'pointer',
  fontSize: '0.9rem',
};

const hiddenInputStyle: React.CSSProperties = {
  display: 'none',
};

const submitButtonStyle: React.CSSProperties = {
  backgroundColor: '#1da1f2',
  color: 'white',
  border: 'none',
  borderRadius: '20px',
  padding: '0.5rem 1.5rem',
  cursor: 'pointer',
  fontWeight: 'bold',
  fontSize: '0.9rem',
};

export default PostForm;