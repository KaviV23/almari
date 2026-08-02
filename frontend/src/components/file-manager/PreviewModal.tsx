"use client"

import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { X, Loader2 } from 'lucide-react'
import { useFileManagerContext } from '@/context/FileManagerContext'

export function PreviewModal() {
  const { 
    showPreviewModal, 
    previewFile, 
    previewContent, 
    previewLoading,
    closePreview,
    isImage,
    isVideo
  } = useFileManagerContext()

  if (!previewFile || !showPreviewModal) return null

  const apiBase = import.meta.env.VITE_API_BASE || ''
  const fileUrl = `${apiBase}/api/download?path=${encodeURIComponent(previewFile.path)}`
  const isMedia = isImage(previewFile.name) || isVideo(previewFile.name)

  return (
    <Dialog open={showPreviewModal} onOpenChange={closePreview}>
      <DialogContent className={isMedia ? "max-w-5xl max-h-[90vh]" : "max-w-4xl max-h-[80vh]"} >
        <DialogHeader className="flex flex-row items-center justify-between border-b border-gray-700 pb-3">
          <DialogTitle className="truncate">{previewFile.name}</DialogTitle>
          <Button variant="ghost" size="icon" onClick={closePreview}>
            <X className="w-5 h-5" />
          </Button>
        </DialogHeader>
        
        {previewLoading ? (
          <div className="flex-1 flex items-center justify-center">
            <Loader2 className="w-8 h-8 animate-spin text-indigo-500" />
          </div>
        ) : isMedia ? (
          <div className="flex-1 flex items-center justify-center overflow-hidden bg-black/60 rounded-xl p-2 min-h-[350px]">
            {isImage(previewFile.name) && (
              <img src={fileUrl} className="max-h-[70vh] max-w-full object-contain rounded-lg" alt={previewFile.name} />
            )}
            {isVideo(previewFile.name) && (
              <video controls preload="metadata" src={fileUrl} className="max-h-[70vh] max-w-full rounded-lg shadow-lg" />
            )}
          </div>
        ) : (
          <Textarea
            value={previewContent}
            className="flex-1 w-full font-mono text-sm min-h-100"
            spellCheck={false}
            readOnly
          />
        )}
        
        <DialogFooter>
          <Button variant="outline" onClick={closePreview}>
            Close
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}