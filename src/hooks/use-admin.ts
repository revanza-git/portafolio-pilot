import { useAuth } from '@/contexts/auth-context';

// Helper function to check if an address is an admin
function isAdminAddress(address: string): boolean {
  const adminAddresses = [
    // Add your admin wallet addresses here
    '0x1234567890123456789012345678901234567890', // Example admin address
  ];
  
  return adminAddresses.includes(address.toLowerCase());
}

export function useAdmin() {
  const { user, isAuthenticated } = useAuth();
  
  const isAdmin = isAuthenticated && user?.address && isAdminAddress(user.address);
  
  return {
    isAdmin: !!isAdmin,
    canAccessAdmin: !!isAdmin,
  };
}