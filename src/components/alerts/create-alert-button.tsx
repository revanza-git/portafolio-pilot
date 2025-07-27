import { Plus } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { CreateAlertModal } from './create-alert-modal';
import { useState } from 'react';

export function CreateAlertButton() {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <>
      <Button onClick={() => setIsOpen(true)} className="bg-gradient-primary hover:opacity-90">
        <Plus className="h-4 w-4 mr-2" />
        Create Alert
      </Button>
      
      <CreateAlertModal 
        isOpen={isOpen} 
        onClose={() => setIsOpen(false)} 
      />
    </>
  );
}