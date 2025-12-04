import { apiClient } from './client'

export interface UploadResult {
  url: string
  filename: string
  size: number
  mimeType: string
}

export type UploadType = 'article' | 'avatar' | 'cover'

export const uploadApi = {
  // Upload image
  uploadImage: async (
    file: File, 
    type: UploadType = 'article',
    onProgress?: (percent: number) => void
  ): Promise<UploadResult> => {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('type', type)
    
    const response = await apiClient.post<UploadResult>('/uploads/image', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (event) => {
        if (event.total && onProgress) {
          const percent = Math.round((event.loaded * 100) / event.total)
          onProgress(percent)
        }
      },
    })
    
    return response.data
  },

  // Delete file
  deleteFile: (url: string) =>
    apiClient.delete('/uploads', { data: { url } }),
}

export default uploadApi


