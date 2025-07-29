import React, { createContext, useContext, useReducer, useEffect } from 'react';
import { toast } from 'sonner';
import { getAPIClient } from '@/lib/api/client';

export interface User {
  id: string;
  address?: string;
  email?: string;
  isEmailVerified?: boolean;
  emailVerified?: boolean; // Add alias for compatibility
  isAdmin?: boolean;
  lastLoginAt?: string; // Add missing property
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

export interface AuthContextType extends AuthState {
  signIn: () => Promise<void>; // Simplified signature to match usage
  signOut: () => void;
  linkEmail: (email: string) => Promise<void>;
  verifyEmail: (code: string) => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

type AuthAction =
  | { type: 'AUTH_START' }
  | { type: 'AUTH_SUCCESS'; payload: User }
  | { type: 'AUTH_ERROR'; payload: string }
  | { type: 'AUTH_LOGOUT' }
  | { type: 'CLEAR_ERROR' };

const initialState: AuthState = {
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null,
};

function authReducer(state: AuthState, action: AuthAction): AuthState {
  switch (action.type) {
    case 'AUTH_START':
      return { ...state, isLoading: true, error: null };
    case 'AUTH_SUCCESS':
      return { ...state, user: action.payload, isAuthenticated: true, isLoading: false, error: null };
    case 'AUTH_ERROR':
      return { ...state, isLoading: false, error: action.payload };
    case 'AUTH_LOGOUT':
      return { ...initialState };
    case 'CLEAR_ERROR':
      return { ...state, error: null };
    default:
      return state;
  }
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [authState, dispatch] = useReducer(authReducer, initialState);

  // Check for existing session on mount
  useEffect(() => {
    const checkExistingSession = async () => {
      console.log('AuthProvider: Checking existing session...');
      const token = localStorage.getItem('auth_token');
      if (token) {
        try {
          dispatch({ type: 'AUTH_START' });
          const apiClient = getAPIClient();
          const userData = await apiClient.getCurrentUser();
          dispatch({ type: 'AUTH_SUCCESS', payload: userData });
          console.log('AuthProvider: Session restored for user:', userData.id);
        } catch (error) {
          console.error('AuthProvider: Session restoration failed:', error);
          localStorage.removeItem('auth_token');
          dispatch({ type: 'AUTH_ERROR', payload: 'Session expired' });
        }
      } else {
        console.log('AuthProvider: No existing session found');
      }
    };

    checkExistingSession();
  }, []);

  const signIn = async () => {
    try {
      console.log('AuthProvider: Starting sign in process...');
      dispatch({ type: 'AUTH_START' });

      // For now, just simulate authentication since SIWE is disabled
      const mockUser: User = {
        id: 'mock-user-id',
        address: '0x1234567890123456789012345678901234567890',
        isEmailVerified: false,
        emailVerified: false,
        isAdmin: false,
        lastLoginAt: new Date().toISOString()
      };

      // Simulate API call
      const apiClient = getAPIClient();
      const authResponse = await apiClient.verifySiwe('mock-message', 'mock-signature');
      
      localStorage.setItem('auth_token', authResponse.token);
      dispatch({ type: 'AUTH_SUCCESS', payload: authResponse.user || mockUser });

      toast.success('Successfully signed in!');
      console.log('AuthProvider: Sign in successful');
    } catch (error) {
      console.error('AuthProvider: Sign in failed:', error);
      dispatch({ type: 'AUTH_ERROR', payload: 'Failed to sign in' });
      toast.error('Failed to sign in. Please try again.');
    }
  };

  const signOut = () => {
    console.log('AuthProvider: Signing out...');
    localStorage.removeItem('auth_token');
    dispatch({ type: 'AUTH_LOGOUT' });
    toast.success('Signed out successfully');
  };

  const linkEmail = async (email: string) => {
    try {
      console.log('AuthProvider: Linking email:', email);
      dispatch({ type: 'AUTH_START' });
      
      // Mock email linking
      toast.success('Verification email sent to ' + email);
      
      dispatch({ type: 'CLEAR_ERROR' });
    } catch (error) {
      console.error('AuthProvider: Email linking failed:', error);
      dispatch({ type: 'AUTH_ERROR', payload: 'Failed to link email' });
      toast.error('Failed to link email. Please try again.');
    }
  };

  const verifyEmail = async (code: string) => {
    try {
      console.log('AuthProvider: Verifying email with code:', code);
      dispatch({ type: 'AUTH_START' });

      // Mock email verification
      if (authState.user) {
        const updatedUser = { ...authState.user, isEmailVerified: true };
        dispatch({ type: 'AUTH_SUCCESS', payload: updatedUser });
      }

      toast.success('Email verified successfully!');
    } catch (error) {
      console.error('AuthProvider: Email verification failed:', error);
      dispatch({ type: 'AUTH_ERROR', payload: 'Failed to verify email' });
      throw error;
    }
  };

  const contextValue: AuthContextType = {
    ...authState,
    signIn,
    signOut,
    linkEmail,
    verifyEmail
  };

  return (
    <AuthContext.Provider value={contextValue}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}