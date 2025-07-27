import { useState } from 'react';
import { Navbar } from '@/components/navigation/navbar';
import { AlertList } from '@/components/alerts/alert-list';
import { CreateAlertButton } from '@/components/alerts/create-alert-button';
import { AlertHistoryComponent } from '@/components/alerts/alert-history';
import { NotificationSettings } from '@/components/alerts/notification-settings';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { useWalletStore } from '@/stores/wallet';
import { Navigate } from 'react-router-dom';

export default function Alerts() {
  const { isConnected } = useWalletStore();

  if (!isConnected) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20">
      <Navbar />
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex justify-between items-center mb-8">
          <div>
            <h1 className="text-3xl font-bold">Price Alerts</h1>
            <p className="text-muted-foreground mt-2">
              Set up notifications for price movements and DeFi events
            </p>
          </div>
          <CreateAlertButton />
        </div>

        <Tabs defaultValue="alerts" className="w-full">
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="alerts">Active Alerts</TabsTrigger>
            <TabsTrigger value="history">History</TabsTrigger>
            <TabsTrigger value="settings">Settings</TabsTrigger>
          </TabsList>
          
          <TabsContent value="alerts" className="space-y-6 mt-6">
            <AlertList />
          </TabsContent>
          
          <TabsContent value="history" className="space-y-6 mt-6">
            <AlertHistoryComponent />
          </TabsContent>
          
          <TabsContent value="settings" className="space-y-6 mt-6">
            <NotificationSettings />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
}