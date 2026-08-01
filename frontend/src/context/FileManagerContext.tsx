/* eslint-disable react/only-export-components */
import React, { createContext, useContext } from 'react'
import { useFileManager } from '@/hooks/useFileManager'

type FileManagerContextType = ReturnType<typeof useFileManager>

const FileManagerContext = createContext<FileManagerContextType | null>(null)

export function FileManagerProvider({ children }: { children: React.ReactNode }) {
  const fileManager = useFileManager()
  return (
    <FileManagerContext.Provider value={fileManager}>
      {children}
    </FileManagerContext.Provider>
  )
}

export function useFileManagerContext() {
  const context = useContext(FileManagerContext)
  if (!context) {
    throw new Error('useFileManagerContext must be used within a FileManagerProvider')
  }
  return context
}
