import axios, { AxiosInstance, AxiosResponse } from 'axios';
import { 
  User, 
  Post, 
  LoginCredentials, 
  RegisterCredentials, 
  CreatePostData, 
  ApiResponse,
  PaginatedResponse 
} from '../types';

class ApiClient {
  private client: AxiosInstance;

  constructor() {
    this.client = axios.create({
      baseURL: process.env.REACT_APP_API_URL || 'http://localhost:8080/api',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    this.client.interceptors.request.use((config) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });

    this.client.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          localStorage.removeItem('token');
          window.location.href = '/login';
        }
        return Promise.reject(error);
      }
    );
  }

  // Auth endpoints
  async login(credentials: LoginCredentials): Promise<ApiResponse<{ user: User; token: string }>> {
    const response: AxiosResponse<ApiResponse<{ user: User; token: string }>> = 
      await this.client.post('/auth/login', credentials);
    return response.data;
  }

  async register(credentials: RegisterCredentials): Promise<ApiResponse<{ user: User; token: string }>> {
    const response: AxiosResponse<ApiResponse<{ user: User; token: string }>> = 
      await this.client.post('/auth/register', credentials);
    return response.data;
  }

  async logout(): Promise<void> {
    await this.client.post('/auth/logout');
  }

  async getCurrentUser(): Promise<ApiResponse<User>> {
    const response: AxiosResponse<ApiResponse<User>> = await this.client.get('/auth/me');
    return response.data;
  }

  // Posts endpoints
  async getPosts(page: number = 1, limit: number = 20): Promise<PaginatedResponse<Post>> {
    const response: AxiosResponse<PaginatedResponse<Post>> = 
      await this.client.get(`/posts?page=${page}&limit=${limit}`);
    return response.data;
  }

  async createPost(data: CreatePostData): Promise<ApiResponse<Post>> {
    const formData = new FormData();
    formData.append('content', data.content);
    
    if (data.imageFiles) {
      data.imageFiles.forEach((file) => {
        formData.append('images', file);
      });
    }

    const response: AxiosResponse<ApiResponse<Post>> = 
      await this.client.post('/posts', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
    return response.data;
  }

  async deletePost(postId: string): Promise<void> {
    await this.client.delete(`/posts/${postId}`);
  }

  async likePost(postId: string): Promise<void> {
    await this.client.post(`/posts/${postId}/like`);
  }

  async unlikePost(postId: string): Promise<void> {
    await this.client.delete(`/posts/${postId}/like`);
  }

  async repost(postId: string): Promise<void> {
    await this.client.post(`/posts/${postId}/repost`);
  }

  async unrepost(postId: string): Promise<void> {
    await this.client.delete(`/posts/${postId}/repost`);
  }

  // Users endpoints
  async getUser(userId: string): Promise<ApiResponse<User>> {
    const response: AxiosResponse<ApiResponse<User>> = await this.client.get(`/users/${userId}`);
    return response.data;
  }

  async getUserPosts(userId: string, page: number = 1, limit: number = 20): Promise<PaginatedResponse<Post>> {
    const response: AxiosResponse<PaginatedResponse<Post>> = 
      await this.client.get(`/users/${userId}/posts?page=${page}&limit=${limit}`);
    return response.data;
  }

  async followUser(userId: string): Promise<void> {
    await this.client.post(`/users/${userId}/follow`);
  }

  async unfollowUser(userId: string): Promise<void> {
    await this.client.delete(`/users/${userId}/follow`);
  }

  async getFollowers(userId: string): Promise<ApiResponse<User[]>> {
    const response: AxiosResponse<ApiResponse<User[]>> = 
      await this.client.get(`/users/${userId}/followers`);
    return response.data;
  }

  async getFollowing(userId: string): Promise<ApiResponse<User[]>> {
    const response: AxiosResponse<ApiResponse<User[]>> = 
      await this.client.get(`/users/${userId}/following`);
    return response.data;
  }
}

export const apiClient = new ApiClient();