"use client"

import { Button } from '@/components/ui/button'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from '@/components/ui/dialog'
import { Checkbox } from '@/components/ui/checkbox'
import { Trash2, X, Loader2 } from 'lucide-react'
import { useFileManagerContext } from '@/context/FileManagerContext'

export function DeleteListModal() {
  const { 
    pendingDeletions, 
    showDeleteListModal, 
    removeFromDeleteList, 
    confirmDeleteAll,
    loading,
    closeDeleteListModal
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
            pendingDeletions.map((path) => (
              <div key={path} className="flex items-center space-x-3 p-3 bg-gray-800/50 border border-gray-700 rounded-lg">
                <Checkbox
                  checked={true}
                  disabled
                  className="text-red-500"
                />
                <span className="flex-1 text-sm font-mono text-gray-300 truncate" title={path}>
                  {path}
                </span>
                <Button
                  variant="ghost"
                  size="icon"
                  className="text-gray-400 hover:text-red-400"
                  onClick={() => removeFromDeleteList(path)}
                  aria-label={`Remove ${path} from delete list`}
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
