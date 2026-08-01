import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { FolderPlus, Upload, HardDrive } from 'lucide-react'
import { useFileManagerContext } from '@/context/FileManagerContext'

export function Header() {
  const { openMkdirModal, uploadFiles } = useFileManagerContext()

  return (
    <header className="bg-gray-800 border-b border-gray-700 px-6 py-4 flex items-center justify-between shadow-md">
      <div className="flex items-center space-x-3">
        <div className="bg-indigo-600 p-2 rounded-lg text-white">
          <HardDrive className="w-6 h-6" />
        </div>
        <div>
          <h1 className="text-xl font-bold tracking-tight text-white">Almari</h1>
          <p className="text-xs text-gray-400">Remote Filesystem Manager</p>
        </div>
      </div>
      <div className="flex items-center space-x-3">
        <Button onClick={openMkdirModal} className="bg-indigo-600 hover:bg-indigo-700">
          <FolderPlus className="w-4 h-4 mr-2" />
          New Folder
        </Button>
        <Label className="flex items-center space-x-2 bg-gray-700 hover:bg-gray-600 text-white px-4 py-2 rounded-lg text-sm font-medium transition cursor-pointer shadow">
          <Upload className="w-4 h-4" />
          <span>Upload Files</span>
          <Input type="file" multiple onChange={uploadFiles} className="hidden" />
        </Label>
      </div>
    </header>
  )
}
