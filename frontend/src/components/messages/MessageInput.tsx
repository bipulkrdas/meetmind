'use client';

import { useState, useRef } from 'react';
import { useParams } from 'next/navigation';
import { messageService } from '@/services/api/message.service';
import { Paperclip, X, Send, Image, FileText, Film } from 'lucide-react';
import { useMessageStore } from '@/stores/message.store';

interface AttachedFile {
  id: string;
  file: File;
  preview?: string;
  uploading: boolean;
  progress: number;
  error?: string;
}

const MAX_FILE_SIZE = 10 * 1024 * 1024; // 10MB
const ALLOWED_FILE_TYPES = {
  images: ['image/jpeg', 'image/png', 'image/gif', 'image/webp'],
  videos: ['video/mp4', 'video/webm', 'video/quicktime'],
  documents: [
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'application/vnd.ms-excel',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'text/plain',
  ],
};

const ALL_ALLOWED_TYPES = [
  ...ALLOWED_FILE_TYPES.images,
  ...ALLOWED_FILE_TYPES.videos,
  ...ALLOWED_FILE_TYPES.documents,
];

export default function MessageInput() {
  const [content, setContent] = useState('');
  const [isSending, setIsSending] = useState(false);
  const [attachedFiles, setAttachedFiles] = useState<AttachedFile[]>([]);
  const [isDragging, setIsDragging] = useState(false);
  const params = useParams();
  const roomId = params.roomId as string;
  const fileInputRef = useRef<HTMLInputElement>(null);
  const dragCounterRef = useRef(0);
  const addMessage = useMessageStore((state) => state.addMessage);

  const validateFile = (file: File): string | null => {
    if (file.size > MAX_FILE_SIZE) {
      return `File "${file.name}" exceeds 10MB limit`;
    }

    if (!ALL_ALLOWED_TYPES.includes(file.type)) {
      return `File type "${file.type}" is not supported`;
    }

    return null;
  };

  const generateFilePreview = (file: File): Promise<string | undefined> => {
    return new Promise((resolve) => {
      if (file.type.startsWith('image/')) {
        const reader = new FileReader();
        reader.onload = (e) => resolve(e.target?.result as string);
        reader.onerror = () => resolve(undefined);
        reader.readAsDataURL(file);
      } else {
        resolve(undefined);
      }
    });
  };

  const handleFileSelect = async (files: FileList | null) => {
    if (!files || files.length === 0) return;

    const newFiles: AttachedFile[] = [];

    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      const error = validateFile(file);

      if (error) {
        alert(error);
        continue;
      }

      const preview = await generateFilePreview(file);

      newFiles.push({
        id: `${Date.now()}-${i}`,
        file,
        preview,
        uploading: false,
        progress: 0,
      });
    }

    setAttachedFiles((prev) => [...prev, ...newFiles]);
  };

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    handleFileSelect(e.target.files);
    e.target.value = '';
  };

  const removeFile = (fileId: string) => {
    setAttachedFiles((prev) => prev.filter((f) => f.id !== fileId));
  };

  const handleDragEnter = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    dragCounterRef.current++;
    if (e.dataTransfer.items && e.dataTransfer.items.length > 0) {
      setIsDragging(true);
    }
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    dragCounterRef.current--;
    if (dragCounterRef.current === 0) {
      setIsDragging(false);
    }
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
  };

  const handleDrop = async (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragging(false);
    dragCounterRef.current = 0;

    const files = e.dataTransfer.files;
    await handleFileSelect(files);
  };

  const uploadAttachedFiles = async (): Promise<string[]> => {
    const uploadedFileIds: string[] = [];
    
    const uploadPromises = attachedFiles.map(async (attachedFile) => {
      try {
        setAttachedFiles((prev) =>
          prev.map((f) =>
            f.id === attachedFile.id ? { ...f, uploading: true, progress: 0 } : f
          )
        );

        const fileId = await messageService.uploadFile(roomId, attachedFile.file, (progress) => {
          setAttachedFiles((prev) =>
            prev.map((f) =>
              f.id === attachedFile.id ? { ...f, progress } : f
            )
          );
        });

        uploadedFileIds.push(fileId);

        setAttachedFiles((prev) =>
          prev.map((f) =>
            f.id === attachedFile.id ? { ...f, uploading: false, progress: 100 } : f
          )
        );
      } catch (error) {
        setAttachedFiles((prev) =>
          prev.map((f) =>
            f.id === attachedFile.id
              ? { ...f, uploading: false, error: 'Upload failed' }
              : f
          )
        );
        throw error;
      }
    });

    await Promise.all(uploadPromises);
    return uploadedFileIds;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if ((!content.trim() && attachedFiles.length === 0) || isSending) return;

    setIsSending(true);

    try {
      let attachment_ids: string[] = [];

      if (attachedFiles.length > 0) {
        attachment_ids = await uploadAttachedFiles();
      }

      const newMessage = await messageService.sendMessage(roomId, {
        content: content.trim(),
        attachment_ids,
      });

      addMessage(newMessage);
      setContent('');
      setAttachedFiles([]);

    } catch (error) {
      console.error('Failed to send message:', error);
      alert('Failed to send message. Please try again.');
    } finally {
      setIsSending(false);
    }
  };

  const getFileIcon = (file: File) => {
    if (file.type.startsWith('image/')) return <Image className="w-5 h-5" />;
    if (file.type.startsWith('video/')) return <Film className="w-5 h-5" />;
    return <FileText className="w-5 h-5" />;
  };

  return (
    <div className="border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800">
      {attachedFiles.length > 0 && (
        <div className="px-4 pt-4 pb-2">
          <div className="flex flex-wrap gap-2">
            {attachedFiles.map((attachedFile) => (
              <div
                key={attachedFile.id}
                className="relative group bg-gray-100 dark:bg-gray-700 rounded-lg p-2 flex items-center gap-2 max-w-xs"
              >
                {attachedFile.preview ? (
                  <img
                    src={attachedFile.preview}
                    alt={attachedFile.file.name}
                    className="w-12 h-12 object-cover rounded"
                  />
                ) : (
                  <div className="w-12 h-12 flex items-center justify-center bg-gray-200 dark:bg-gray-600 rounded text-gray-600 dark:text-gray-300">
                    {getFileIcon(attachedFile.file)}
                  </div>
                )}

                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">
                    {attachedFile.file.name}
                  </p>
                  <p className="text-xs text-gray-500 dark:text-gray-400">
                    {(attachedFile.file.size / 1024).toFixed(1)} KB
                  </p>

                  {attachedFile.uploading && (
                    <div className="mt-1">
                      <div className="w-full bg-gray-200 dark:bg-gray-600 rounded-full h-1">
                        <div
                          className="bg-blue-600 h-1 rounded-full transition-all duration-300"
                          style={{ width: `${attachedFile.progress}%` }}
                        />
                      </div>
                    </div>
                  )}

                  {attachedFile.error && (
                    <p className="text-xs text-red-600 dark:text-red-400 mt-1">
                      {attachedFile.error}
                    </p>
                  )}
                </div>

                {!attachedFile.uploading && (
                  <button
                    onClick={() => removeFile(attachedFile.id)}
                    className="absolute -top-2 -right-2 bg-red-500 hover:bg-red-600 text-white rounded-full p-1 opacity-0 group-hover:opacity-100 transition-opacity"
                  >
                    <X className="w-3 h-3" />
                  </button>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      <div className="p-4">
        <form
          onSubmit={handleSubmit}
          onDragEnter={handleDragEnter}
          onDragLeave={handleDragLeave}
          onDragOver={handleDragOver}
          onDrop={handleDrop}
          className={`relative flex gap-2 ${
            isDragging
              ? 'ring-2 ring-blue-500 ring-offset-2 dark:ring-offset-gray-800 rounded-lg'
              : ''
          }`}
        >
          {isDragging && (
            <div className="absolute inset-0 bg-blue-500/10 dark:bg-blue-500/20 border-2 border-dashed border-blue-500 rounded-lg flex items-center justify-center z-10 pointer-events-none">
              <div className="text-center">
                <Paperclip className="w-12 h-12 mx-auto mb-2 text-blue-600 dark:text-blue-400" />
                <p className="text-sm font-medium text-blue-600 dark:text-blue-400">
                  Drop files to attach
                </p>
              </div>
            </div>
          )}

          <input
            ref={fileInputRef}
            type="file"
            multiple
            accept={ALL_ALLOWED_TYPES.join(',')}
            onChange={handleFileInputChange}
            className="hidden"
          />

          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            disabled={isSending}
            className="flex-shrink-0 p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            title="Attach files"
          >
            <Paperclip className="w-5 h-5" />
          </button>

          <input
            type="text"
            value={content}
            onChange={(e) => {
              setContent(e.target.value);
            }}
            placeholder="Type a message..."
            className="flex-1 px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white dark:bg-gray-900 text-gray-900 dark:text-gray-100 placeholder-gray-500 dark:placeholder-gray-400"
            disabled={isSending}
          />

          <button
            type="submit"
            disabled={(!content.trim() && attachedFiles.length === 0) || isSending}
            className="flex-shrink-0 px-4 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 dark:disabled:bg-gray-600 text-white rounded-lg font-medium transition-colors disabled:cursor-not-allowed flex items-center gap-2"
          >
            {isSending ? (
              <>
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white" />
                Sending...
              </>
            ) : (
              <>
                <Send className="w-4 h-4" />
                Send
              </>
            )}
          </button>
        </form>

        <p className="text-xs text-gray-500 dark:text-gray-400 mt-2">
          Supported: Images, Videos, Documents (PDF, Word, Excel) • Max 10MB per file
        </p>
      </div>
    </div>
  );
}
