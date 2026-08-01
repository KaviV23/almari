import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useFileManagerContext } from '@/context/FileManagerContext'

export function MkdirModal() {
  const { showMkdirModal, newFolderName, setNewFolderName, createFolder, closeMkdirModal } = useFileManagerContext()

  return (
    <Dialog open={showMkdirModal} onOpenChange={closeMkdirModal}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>Create New Folder</DialogTitle>
        </DialogHeader>
        <div className="space-y-4">
          <div className="space-y-2">
            <div><Label htmlFor="folder-name">Folder name</Label></div>
            <Input
              id="folder-name"
              placeholder="Folder name"
              value={newFolderName}
              onChange={(e) => setNewFolderName(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && createFolder()}
              autoFocus
            />
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={closeMkdirModal}>
            Cancel
          </Button>
          <Button onClick={createFolder} disabled={!newFolderName}>
            Create
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
