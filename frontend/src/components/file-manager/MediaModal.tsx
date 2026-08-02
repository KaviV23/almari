import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { X } from 'lucide-react'
import { useFileManagerContext } from '@/context/FileManagerContext'

export function MediaModal() {
  const { showMediaModal, currentFile, closeMediaModal, isImage, isVideo } = useFileManagerContext()

  if (!currentFile || !showMediaModal) return null

  const apiBase = import.meta.env.VITE_API_BASE || ''
  const fileUrl = `${apiBase}/api/download?path=${encodeURIComponent(currentFile.path)}`

  // Only show for images/videos
  if (!isImage(currentFile.name) && !isVideo(currentFile.name)) return null

  return (
    <Dialog open={showMediaModal} onOpenChange={closeMediaModal}>
      <DialogContent className="max-w-5xl max-h-[90vh]">
        <DialogHeader className="flex flex-row items-center justify-between border-b border-gray-700 pb-3">
          <DialogTitle className="truncate">{currentFile.name}</DialogTitle>
          <Button variant="ghost" size="icon" onClick={closeMediaModal}>
            <X className="w-5 h-5" />
          </Button>
        </DialogHeader>
        <div className="flex-1 flex items-center justify-center overflow-hidden bg-black/60 rounded-xl p-2 min-h-[350px]">
          {isImage(currentFile.name) && (
            <img src={fileUrl} className="max-h-[70vh] max-w-full object-contain rounded-lg" alt={currentFile.name} />
          )}
          {isVideo(currentFile.name) && (
            <video controls preload="metadata" src={fileUrl} className="max-h-[70vh] max-w-full rounded-lg shadow-lg" />
          )}
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={closeMediaModal}>
            Close
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}