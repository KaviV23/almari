export interface FileItem {
  name: string
  path: string
  size: number
  mode: string
  modTime: string
  isDir: boolean
  hasThumb: boolean
}

export interface FileManagerState {
  currentPath: string
  files: FileItem[]
  breadcrumbs: string[]
  loading: boolean
  viewMode: 'list' | 'icon'
  showMkdirModal: boolean
  newFolderName: string
  showEditorModal: boolean
  showMediaModal: boolean
  showDeleteListModal: boolean
  currentFile: FileItem | null
  fileContent: string
  selectedFiles: string[]
  pendingDeletions: string[]
}