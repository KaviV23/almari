"use client"

import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog'
import { Trash2, X, Loader2, Folder, File, Image, Video } from 'lucide-react'
import { useFileManagerContext } from '@/context/FileManagerContext'

export function DeleteListModal() {
  const { 
    pendingDeletions, 
    showDeleteListModal, 
    removeFromDeleteList, 
    confirmDeleteAll,
    loading,
    closeDeleteListModal,
    isImage,
    isVideo
  } = useFileManagerContext()

  if (!showDeleteListModal) {
    return null
  }

  return (
    <Dialog open={showDeleteListModal} onOpenChange={closeDeleteListModal}>
      <DialogContent className="max-w-2xl">
        <DialogHeader>
          <DialogTitle className="flex items-center space-x-2">
            <Trash2 className="w-5 h-5 text-red-400" />
            <span>Delete List ({pendingDeletions.length})</span>
          </DialogTitle>
          <DialogDescription>
            Review the items below before confirming deletion. You can remove items from this list.
          </DialogDescription>
        </DialogHeader>
        <div className="max-h-96 overflow-y-auto space-y-2">
          {pendingDeletions.length === 0 ? (
            <div className="p-4 text-center text-gray-400">
              No items in delete list
            </div>
          ) : (
            pendingDeletions.map((file) => (
              <div key={file.path} className="flex items-center space-x-3 p-3 bg-gray-800/50 border border-gray-700 rounded-lg">
                <div className="w-12 h-12 shrink-0 bg-gray-900 rounded flex items-center justify-center overflow-hidden relative">
                  {file.hasThumb ? (
                    <img
                      src={`${import.meta.env.VITE_API_BASE || ''}/api/thumbnail?path=${encodeURIComponent(file.path)}`}
                      className="w-full h-full object-cover"
                      alt={file.name}
                    />
                  ) : file.isDir ? (
                    <Folder className="w-6 h-6 text-amber-400" />
                  ) : isImage(file.name) ? (
                    <Image className="w-6 h-6 text-indigo-400" />
                  ) : isVideo(file.name) ? (
                    <Video className="w-6 h-6 text-indigo-400" />
                  ) : (
                    <File className="w-6 h-6 text-indigo-400" />
                  )}
                </div>
                <div className="flex-1 min-w-0">
                  <p className="text-sm font-medium text-gray-200 truncate" title={file.name}>
                    {file.name}
                  </p>
                  <p className="text-xs font-mono text-gray-500 truncate" title={file.path}>
                    {file.path}
                  </p>
                </div>
                <Button
                  variant="ghost"
                  size="icon"
                  className="text-gray-400 hover:text-red-400"
                  onClick={() => removeFromDeleteList(file.path)}
                  aria-label={`Remove ${file.name} from delete list`}
                >
                  <X className="w-4 h-4" />
                </Button>
              </div>
            ))
          )}
        </div>
        <DialogFooter className="flex flex-row space-x-2">
          <Button variant="outline" onClick={closeDeleteListModal} disabled={loading}>
            Cancel
          </Button>
          <Button variant="destructive" onClick={confirmDeleteAll} disabled={loading || pendingDeletions.length === 0}>
            {loading ? (
              <>
                <Loader2 className="w-4 h-4 animate-spin mr-2" />
                Deleting...
              </>
            ) : (
              <>
                <Trash2 className="w-4 h-4 mr-2" />
                Confirm Delete All
              </>
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
