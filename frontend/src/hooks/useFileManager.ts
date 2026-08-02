import { useState, useCallback, useEffect } from 'react'
import type { FileItem, FileManagerState } from '@/types/file'

const API_BASE = ''

function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function formatDate(timestamp: string): string {
  return new Date(timestamp).toLocaleString()
}

function isImage(filename: string): boolean {
  const name = filename.toLowerCase()
  return name.endsWith('.jpg') || name.endsWith('.jpeg') || name.endsWith('.png') || name.endsWith('.gif') || name.endsWith('.webp')
}

function isVideo(filename: string): boolean {
  const name = filename.toLowerCase()
  return name.endsWith('.mp4') || name.endsWith('.mkv') || name.endsWith('.mov') || name.endsWith('.webm') || name.endsWith('.avi')
}

export function useFileManager() {
  const [state, setState] = useState<FileManagerState>({
    currentPath: '/',
    files: [],
    breadcrumbs: [],
    loading: false,
    viewMode: 'list',
    showMkdirModal: false,
    newFolderName: '',
    showEditorModal: false,
    showMediaModal: false,
    showDeleteListModal: false,
    showPreviewModal: false,
    currentFile: null,
    fileContent: '',
    selectedFiles: [],
    pendingDeletions: [],
    previewFile: null,
    previewContent: '',
    previewLoading: false,
  })

  const updateBreadcrumbs = useCallback((path: string) => {
    if (path === '/' || path === '') {
      setState(prev => ({ ...prev, breadcrumbs: [] }))
    } else {
      setState(prev => ({ ...prev, breadcrumbs: path.split('/').filter(Boolean) }))
    }
  }, [])

  const loadDir = useCallback(async (path: string) => {
    setState(prev => ({ ...prev, loading: true, currentPath: path }))
    updateBreadcrumbs(path)

    try {
      const res = await fetch(`${API_BASE}/api/list?path=${encodeURIComponent(path)}`)
      if (res.ok) {
        const files = await res.json()
        setState(prev => ({ ...prev, files, loading: false }))
      } else {
        alert('Failed to load directory')
        setState(prev => ({ ...prev, loading: false }))
      }
    } catch (e) {
      console.error(e)
      alert('Network error')
      setState(prev => ({ ...prev, loading: false }))
    }
  }, [updateBreadcrumbs])

  const navigateToBreadcrumb = useCallback((index: number) => {
    const targetPath = '/' + state.breadcrumbs.slice(0, index + 1).join('/')
    loadDir(targetPath)
  }, [state.breadcrumbs, loadDir])

  const openMkdirModal = useCallback(() => {
    setState(prev => ({ ...prev, newFolderName: '', showMkdirModal: true }))
  }, [])

  const createFolder = useCallback(async () => {
    const { newFolderName, currentPath } = state
    if (!newFolderName) return

    const target = (currentPath === '/' ? '' : currentPath) + '/' + newFolderName

    try {
      const res = await fetch(`${API_BASE}/api/mkdir`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ path: target })
      })
      if (res.ok) {
        setState(prev => ({ ...prev, showMkdirModal: false }))
        loadDir(currentPath)
      } else {
        alert('Failed to create folder')
      }
    } catch (e) {
      console.error(e)
    }
  }, [state.newFolderName, state.currentPath, loadDir])

  const deleteItem = useCallback(async (path: string) => {
    if (!confirm(`Are you sure you want to delete ${path}?`)) return

    try {
      const res = await fetch(`${API_BASE}/api/delete?path=${encodeURIComponent(path)}`, {
        method: 'DELETE'
      })
      if (res.ok) {
        loadDir(state.currentPath)
      } else {
        alert('Failed to delete item')
      }
    } catch (e) {
      console.error(e)
    }
  }, [state.currentPath, loadDir])

  const toggleSelectFile = useCallback((path: string) => {
    setState(prev => ({
      ...prev,
      selectedFiles: prev.selectedFiles.includes(path)
        ? prev.selectedFiles.filter(p => p !== path)
        : [...prev.selectedFiles, path]
    }))
  }, [])

  const selectAllFiles = useCallback(() => {
    setState(prev => ({
      ...prev,
      selectedFiles: prev.selectedFiles.length === prev.files.length
        ? []
        : prev.files.map(f => f.path)
    }))
  }, [])

  const clearSelection = useCallback(() => {
    setState(prev => ({ ...prev, selectedFiles: [] }))
  }, [])

  const addToDeleteList = useCallback(() => {
    setState(prev => {
      const selectedFileItems = prev.files.filter(f => prev.selectedFiles.includes(f.path))
      const existingPaths = new Set(prev.pendingDeletions.map(f => f.path))
      const newItems = selectedFileItems.filter(f => !existingPaths.has(f.path))
      return {
        ...prev,
        pendingDeletions: [...prev.pendingDeletions, ...newItems],
        selectedFiles: [],
      }
    })
  }, [])

  const removeFromDeleteList = useCallback((path: string) => {
    setState(prev => ({
      ...prev,
      pendingDeletions: prev.pendingDeletions.filter(f => f.path !== path)
    }))
  }, [])

  const clearDeleteList = useCallback(() => {
    setState(prev => ({ ...prev, pendingDeletions: [] }))
  }, [])

  const confirmDeleteAll = useCallback(async () => {
    if (state.pendingDeletions.length === 0) return

    try {
      for (const file of state.pendingDeletions) {
        const res = await fetch(`${API_BASE}/api/delete?path=${encodeURIComponent(file.path)}`, {
          method: 'DELETE'
        })
        if (!res.ok) {
          alert(`Failed to delete ${file.path}`)
          return
        }
      }
      setState(prev => ({ ...prev, pendingDeletions: [] }))
      loadDir(state.currentPath)
    } catch (e) {
      console.error(e)
      alert('Network error during deletion')
    }
  }, [state.pendingDeletions, state.currentPath, loadDir])

  const openDeleteListModal = useCallback(() => {
    setState(prev => ({ ...prev, showDeleteListModal: true }))
  }, [])

  const closeDeleteListModal = useCallback(() => {
    setState(prev => ({ ...prev, showDeleteListModal: false }))
  }, [])

  const openPreview = useCallback(async (file: FileItem) => {
    if (file.isDir) return
    setState(prev => ({ ...prev, previewLoading: true, previewFile: file }))
    try {
      const res = await fetch(`${API_BASE}/api/read?path=${encodeURIComponent(file.path)}`)
      if (res.ok) {
        const content = await res.text()
        setState(prev => ({ ...prev, previewContent: content, previewLoading: false, showPreviewModal: true }))
      } else {
        alert('Failed to load preview')
        setState(prev => ({ ...prev, previewLoading: false, previewFile: null }))
      }
    } catch (e) {
      console.error(e)
      alert('Network error')
      setState(prev => ({ ...prev, previewLoading: false, previewFile: null }))
    }
  }, [])

  const closePreview = useCallback(() => {
    setState(prev => ({ ...prev, previewFile: null, previewContent: '', showPreviewModal: false }))
  }, [])

  const uploadFiles = useCallback(async (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files
    if (!files?.length) return

    const formData = new FormData()
    formData.append('path', state.currentPath)
    for (let i = 0; i < files.length; i++) {
      formData.append('files', files[i])
    }

    try {
      const res = await fetch(`${API_BASE}/api/upload`, {
        method: 'POST',
        body: formData
      })
      if (res.ok) {
        loadDir(state.currentPath)
      } else {
        alert('Failed to upload files')
      }
    } catch (e) {
      console.error(e)
    }
    event.target.value = ''
  }, [state.currentPath, loadDir])

  const loadTextFile = useCallback(async (file: FileItem) => {
    try {
      const res = await fetch(`${API_BASE}/api/read?path=${encodeURIComponent(file.path)}`)
      if (res.ok) {
        const content = await res.text()
        setState(prev => ({ ...prev, fileContent: content, showEditorModal: true }))
      } else {
        alert('Failed to read file')
      }
    } catch (e) {
      console.error(e)
    }
  }, [])

  const openFile = useCallback((file: FileItem) => {
    setState(prev => ({ ...prev, currentFile: file }))
    const isImg = isImage(file.name)
    const isVid = isVideo(file.name)

    if (isImg || isVid) {
      setState(prev => ({ ...prev, showMediaModal: true }))
    } else {
      loadTextFile(file)
    }
  }, [loadTextFile])

  const closeMediaModal = useCallback(() => {
    setState(prev => ({ ...prev, showMediaModal: false, currentFile: null }))
  }, [])

  const saveFile = useCallback(async () => {
    const { currentFile, fileContent } = state
    if (!currentFile) return

    try {
      const res = await fetch(`${API_BASE}/api/write`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ path: currentFile.path, content: fileContent })
      })
      if (res.ok) {
        alert('File saved successfully')
        setState(prev => ({ ...prev, showEditorModal: false }))
      } else {
        alert('Failed to save file')
      }
    } catch (e) {
      console.error(e)
    }
  }, [state.currentFile, state.fileContent])

  // Load root directory on mount
  useEffect(() => {
    loadDir('/')
  }, [loadDir])

  return {
    ...state,
    formatSize,
    formatDate,
    isImage,
    isVideo,
    loadDir,
    navigateToBreadcrumb,
    openMkdirModal,
    createFolder,
    deleteItem,
    uploadFiles,
    openFile,
    closeMediaModal,
    saveFile,
    setViewMode: (mode: 'list' | 'icon') => setState(prev => ({ ...prev, viewMode: mode })),
    setNewFolderName: (name: string) => setState(prev => ({ ...prev, newFolderName: name })),
    setFileContent: (content: string) => setState(prev => ({ ...prev, fileContent: content })),
    closeMkdirModal: () => setState(prev => ({ ...prev, showMkdirModal: false })),
    closeEditorModal: () => setState(prev => ({ ...prev, showEditorModal: false, currentFile: null })),
    toggleSelectFile,
    selectAllFiles,
    clearSelection,
    addToDeleteList,
    removeFromDeleteList,
    clearDeleteList,
    confirmDeleteAll,
    openDeleteListModal,
    closeDeleteListModal,
    openPreview,
    closePreview,
  }
}