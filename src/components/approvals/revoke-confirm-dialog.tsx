import { AlertTriangle, Loader2 } from 'lucide-react';
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
import { useRevokeAllowance } from '@/hooks/use-revoke-allowance';
import { Address } from 'viem';

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
  const { revokeAllowance, isLoading, isSuccess, error } = useRevokeAllowance();

  const handleConfirm = async () => {
    if (!allowance) return;
    
    await revokeAllowance({
      tokenAddress: allowance.token.address as Address,
      spenderAddress: allowance.spender.address as Address,
      tokenSymbol: allowance.token.symbol,
      spenderName: allowance.spender.name,
    });
    
    if (!error) {
      onConfirm();
      onClose();
    }
  };

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
          <AlertDialogCancel onClick={onClose} disabled={isLoading}>
            Cancel
          </AlertDialogCancel>
          <AlertDialogAction 
            onClick={handleConfirm}
            disabled={isLoading}
            className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
          >
            {isLoading ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                Revoking...
              </>
            ) : (
              'Revoke Allowance'
            )}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}