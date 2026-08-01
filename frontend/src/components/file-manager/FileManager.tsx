import { Header } from './Header'
import { Breadcrumbs } from './Breadcrumbs'
import { FileList } from './FileList'
import { FileGrid } from './FileGrid'
import { MkdirModal } from './MkdirModal'
import { EditorModal } from './EditorModal'
import { MediaModal } from './MediaModal'
import { FileManagerProvider, useFileManagerContext } from '@/context/FileManagerContext'

function FileManagerContent() {
  const { viewMode } = useFileManagerContext()

  return (
    <div className="min-h-screen bg-background text-foreground flex flex-col">
      <Header />
      <main className="flex-1 max-w-7xl w-full mx-auto p-6 flex flex-col space-y-4">
        <Breadcrumbs />
        <div className="bg-card border border-border rounded-xl shadow overflow-hidden flex-1 flex flex-col">
          {viewMode === 'list' ? (
            <div className="overflow-x-auto flex-1">
              <FileList />
            </div>
          ) : (
            <FileGrid />
          )}
        </div>
      </main>
      <MkdirModal />
      <EditorModal />
      <MediaModal />
    </div>
  )
}

export function FileManager() {
  return (
    <FileManagerProvider>
      <FileManagerContent />
    </FileManagerProvider>
  )
}
