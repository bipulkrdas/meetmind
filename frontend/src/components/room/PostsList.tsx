'use client';

import { useEffect, useState } from 'react';
import { postService, CreatePostRequest } from '@/services/api/post.service';

interface PostsListProps {
  roomId: string;
}

export default function PostsList({ roomId }: PostsListProps) {
  const [posts, setPosts] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [newMessage, setNewMessage] = useState('');

  useEffect(() => {
    loadPosts();
  }, [roomId]);

  const loadPosts = async () => {
    try {
      const data = await postService.getPosts(roomId);
      setPosts(data);
    } catch (error) {
      console.error('Failed to load posts:', error);
    } finally {
      setLoading(false);
    }
  };

  const handlePostSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim()) return;

    try {
      await postService.createPost(roomId, { message: newMessage });
      setNewMessage('');
      loadPosts(); // Refresh posts
    } catch (error) {
      console.error('Failed to create post:', error);
    }
  };

  return (
    <div className="mt-8">
      <h2 className="text-lg font-bold mb-4">Posts</h2>
      <form onSubmit={handlePostSubmit} className="mb-4">
        <textarea
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="Write a post..."
          className="w-full p-2 border rounded"
          rows={3}
        />
        <button type="submit" className="mt-2 px-4 py-2 bg-blue-600 text-white rounded">
          Post
        </button>
      </form>
      {loading ? (
        <div>Loading posts...</div>
      ) : (
        <div className="space-y-4">
          {posts.map(post => (
            <div key={post.id} className="p-4 bg-white rounded-lg shadow-md">
              <div className="flex items-center mb-2">
                <span className="font-bold">{post.creator_name}</span>
                <span className="text-sm text-gray-500 ml-2">{new Date(post.created_at).toLocaleString()}</span>
              </div>
              <p>{post.message}</p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
