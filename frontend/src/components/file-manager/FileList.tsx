import { Table, TableHeader, TableBody, TableRow, TableHead, TableCell } from '@/components/ui/table'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipTrigger, TooltipContent, TooltipProvider } from '@/components/ui/tooltip'
import { Checkbox } from '@/components/ui/checkbox'
import { Folder, File, Download, Trash2, Loader2 } from 'lucide-react'
import { useFileManagerContext } from '@/context/FileManagerContext'
import type { FileItem } from '@/types/file'

export function FileList() {
  const { files, loading, loadDir, openFile, deleteItem, formatSize, formatDate, selectedFiles, toggleSelectFile, selectAllFiles, addToDeleteList } = useFileManagerContext()

  if (loading) {
    return (
      <TooltipProvider>
        <div className="flex justify-center items-center h-full p-10">
          <Loader2 className="w-5 h-5 animate-spin text-primary mr-2" />
          <span>Loading directory contents...</span>
        </div>
      </TooltipProvider>
    )
  }

  if (files.length === 0) {
    return (
      <div className="flex justify-center items-center h-full text-gray-500">
        This directory is empty.
      </div>
    )
  }

  return (
    <TooltipProvider>
      <Table>
        <TableHeader>
          <TableRow className="text-xs font-semibold uppercase tracking-wider">
            <TableHead className="py-3 px-6 w-12">
              <Checkbox
                checked={selectedFiles.length === files.length && files.length > 0}
                onCheckedChange={selectAllFiles}
                aria-label="Select all files"
              />
            </TableHead>
            <TableHead className="py-3 px-6">Name</TableHead>
            <TableHead className="py-3 px-6">Size</TableHead>
            <TableHead className="py-3 px-6">Permissions</TableHead>
            <TableHead className="py-3 px-6">Modified</TableHead>
            <TableHead className="py-3 px-6 text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody className="divide-y divide-gray-700 text-sm">
          {files.map((file: FileItem) => (
            <TableRow key={file.path} className="hover:bg-gray-700/50 transition">
              <TableCell className="py-3 px-6 w-12">
                <Checkbox
                  checked={selectedFiles.includes(file.path)}
                  onCheckedChange={() => toggleSelectFile(file.path)}
                  aria-label={`Select ${file.name}`}
                />
              </TableCell>
              <TableCell className="py-3 px-6">
                <div className="flex items-center space-x-3">
                  {file.isDir ? (
                    <Folder className="w-5 h-5 text-amber-400 shrink-0" />
                  ) : (
                    <File className="w-5 h-5 text-indigo-400 shrink-0" />
                  )}
                  <Button
                    variant="ghost"
                    size="sm"
                    className="font-medium text-gray-200 hover:text-indigo-400 text-left truncate max-w-xs h-auto p-0"
                    onClick={() => file.isDir ? loadDir(file.path) : openFile(file)}
                  >
                    {file.name}
                  </Button>
                </div>
              </TableCell>
              <TableCell className="py-3 px-6 text-gray-400 text-xs">
                {file.isDir ? '-' : formatSize(file.size)}
              </TableCell>
              <TableCell className="py-3 px-6 text-gray-400 font-mono text-xs">
                {file.mode}
              </TableCell>
              <TableCell className="py-3 px-6 text-gray-400 text-xs">
                {formatDate(file.modTime)}
              </TableCell>
              <TableCell className="py-3 px-6">
                <div className="flex items-center justify-end space-x-2">
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <a
                        href={`${import.meta.env.VITE_API_BASE || ''}/api/download?path=${encodeURIComponent(file.path)}`}
                        download
                        className="text-gray-400 hover:text-indigo-400 p-1 rounded transition inline-flex items-center"
                      >
                        <Download className="w-4 h-4" />
                      </a>
                    </TooltipTrigger>
                    <TooltipContent>Download</TooltipContent>
                  </Tooltip>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <Button
                        variant="ghost"
                        size="icon"
                        className="text-gray-400 hover:text-red-400"
                        onClick={() => deleteItem(file.path)}
                      >
                        <Trash2 className="w-4 h-4" />
                      </Button>
                    </TooltipTrigger>
                    <TooltipContent>Delete</TooltipContent>
                  </Tooltip>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      {selectedFiles.length > 0 && (
        <div className="fixed bottom-4 right-4 bg-card border border-border rounded-lg shadow-lg p-4 flex items-center space-x-2 z-50">
          <span className="text-sm text-gray-300">{selectedFiles.length} selected</span>
          <Button size="sm" variant="outline" onClick={addToDeleteList}>
            Add to Delete List
          </Button>
        </div>
      )}
    </TooltipProvider>
  )
}
