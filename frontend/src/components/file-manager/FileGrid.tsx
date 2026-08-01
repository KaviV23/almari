import { Button } from '@/components/ui/button'
import { Tooltip, TooltipTrigger, TooltipContent, TooltipProvider } from '@/components/ui/tooltip'
import { Folder, File, Download, Trash2, Loader2, Image, Video } from 'lucide-react'
import { useFileManagerContext } from '@/context/FileManagerContext'
import type { FileItem } from '@/types/file'

export function FileGrid() {
  const { files, loading, loadDir, openFile, deleteItem, isImage, isVideo } = useFileManagerContext()

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
      <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4 p-6 h-full overflow-y-auto">
        {files.map((file: FileItem) => (
          <div
            key={file.path}
            className="bg-card border border-border rounded p-2 flex flex-col items-center space-y-2 hover:border-indigo-500 transition relative group"
          >
            <div className="w-full h-28 bg-gray-800 rounded flex items-center justify-center overflow-hidden relative">
              {file.hasThumb ? (
                <img
                  src={`${import.meta.env.VITE_API_BASE || ''}/api/thumbnail?path=${encodeURIComponent(file.path)}`}
                  className="w-full h-full object-cover"
                  alt={file.name}
                />
              ) : file.isDir ? (
                <Folder className="w-10 h-10 text-amber-400" />
              ) : isImage(file.name) ? (
                <Image className="w-10 h-10 text-indigo-400" />
              ) : isVideo(file.name) ? (
                <Video className="w-10 h-10 text-indigo-400" />
              ) : (
                <File className="w-10 h-10 text-indigo-400" />
              )}
            </div>
            <Button
              variant="ghost"
              size="sm"
              className="text-xs font-medium text-gray-200 hover:text-indigo-400 text-center truncate w-full h-auto p-0"
              onClick={() => file.isDir ? loadDir(file.path) : openFile(file)}
              title={file.name}
            >
              {file.name}
            </Button>
            <div className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition flex space-x-1 bg-gray-800/90 backdrop-blur rounded p-1 shadow">
              <Tooltip>
                <TooltipTrigger asChild>
                  <a
                    href={`${import.meta.env.VITE_API_BASE || ''}/api/download?path=${encodeURIComponent(file.path)}`}
                    download
                    className="text-gray-400 hover:text-indigo-400 p-1"
                  >
                    <Download className="w-3.5 h-3.5" />
                  </a>
                </TooltipTrigger>
                <TooltipContent>Download</TooltipContent>
              </Tooltip>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button
                    variant="ghost"
                    size="icon"
                    className="text-gray-400 hover:text-red-400 h-6 w-6 p-0"
                    onClick={() => deleteItem(file.path)}
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </Button>
                </TooltipTrigger>
                <TooltipContent>Delete</TooltipContent>
              </Tooltip>
            </div>
          </div>
        ))}
      </div>
    </TooltipProvider>
  )
}
