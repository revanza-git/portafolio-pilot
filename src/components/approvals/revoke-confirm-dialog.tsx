import { AlertTriangle } from 'lucide-react';
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from '@/components/ui/alert-dialog';
import { Allowance } from '@/stores/portfolio';

interface RevokeConfirmDialogProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  allowance?: Allowance;
}

export function RevokeConfirmDialog({ 
  isOpen, 
  onClose, 
  onConfirm, 
  allowance 
}: RevokeConfirmDialogProps) {
  if (!allowance) return null;

  return (
    <AlertDialog open={isOpen} onOpenChange={onClose}>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle className="flex items-center gap-2">
            <AlertTriangle className="h-5 w-5 text-warning" />
            Revoke Token Allowance
          </AlertDialogTitle>
          <AlertDialogDescription className="space-y-2">
            <p>
              You are about to revoke the allowance for{' '}
              <strong>{allowance.spender.name}</strong> to spend your{' '}
              <strong>{allowance.token.symbol}</strong> tokens.
            </p>
            <div className="bg-muted p-3 rounded-lg space-y-1">
              <div className="text-sm">
                <strong>Token:</strong> {allowance.token.name} ({allowance.token.symbol})
              </div>
              <div className="text-sm">
                <strong>Spender:</strong> {allowance.spender.name}
              </div>
              <div className="text-sm">
                <strong>Current Allowance:</strong> {allowance.amountFormatted}
              </div>
            </div>
            <p className="text-xs text-muted-foreground">
              This action requires a transaction and will cost gas fees. 
              After revoking, you'll need to approve again if you want to use this protocol.
            </p>
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel onClick={onClose}>Cancel</AlertDialogCancel>
          <AlertDialogAction 
            onClick={onConfirm}
            className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
          >
            Revoke Allowance
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}