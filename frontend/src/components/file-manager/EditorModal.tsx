import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { X } from 'lucide-react'
import { useFileManagerContext } from '@/context/FileManagerContext'

export function EditorModal() {
  const { showEditorModal, currentFile, fileContent, setFileContent, saveFile, closeEditorModal } = useFileManagerContext()

  if (!currentFile) return null

  return (
    <Dialog open={showEditorModal} onOpenChange={closeEditorModal}>
      <DialogContent className="max-w-4xl max-h-[80vh]">
        <DialogHeader className="flex flex-row items-center justify-between border-b border-gray-700 pb-3">
          <DialogTitle className="truncate">{currentFile.name}</DialogTitle>
          <Button variant="ghost" size="icon" onClick={closeEditorModal}>
            <X className="w-5 h-5" />
          </Button>
        </DialogHeader>
        <Textarea
          value={fileContent}
          onChange={(e) => setFileContent(e.target.value)}
          className="flex-1 w-full font-mono text-sm min-h-100"
          spellCheck={false}
        />
        <DialogFooter>
          <Button variant="outline" onClick={closeEditorModal}>
            Close
          </Button>
          <Button onClick={saveFile}>
            Save Changes
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
