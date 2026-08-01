import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { FolderPlus, Upload, HardDrive } from 'lucide-react'
import { useFileManagerContext } from '@/context/FileManagerContext'

export function Header() {
  const { openMkdirModal, uploadFiles } = useFileManagerContext()

  return (
    <header className="bg-background border-b border-border p-2 flex items-center justify-between shadow-md">
      <div className="flex items-center space-x-3">
        <div className="bg-accent p-2 rounded-lg text-accent-foreground">
          <HardDrive className="w-6 h-6" />
        </div>
        <div>
          <h1 className="text-md font-bold tracking-tight text-foreground">Almari</h1>
          <p className="text-xs text-muted-foreground">Remote Filesystem Manager</p>
        </div>
      </div>
      <div className="flex items-center space-x-3">
        <Button onClick={openMkdirModal} size={"sm"}>
          <FolderPlus className="w-4 h-4" />
          New Folder
        </Button>
        <Label className="flex items-center space-x-2 bg-secondary hover:bg-border text-secondary-foreground px-3 h-9 rounded-md text-sm font-medium transition cursor-pointer shadow">
          <Upload className="w-4 h-4" />
          <span>Upload Files</span>
          <Input type="file" multiple onChange={uploadFiles} className="hidden" />
        </Label>
      </div>
    </header>
  )
}
