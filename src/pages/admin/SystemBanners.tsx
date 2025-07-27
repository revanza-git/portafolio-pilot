import { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Switch } from '@/components/ui/switch';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { 
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { 
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { 
  Plus, 
  Edit, 
  Megaphone, 
  RefreshCw,
  AlertTriangle,
  Info,
  CheckCircle,
  AlertCircle,
  Trash2
} from 'lucide-react';
import { useFeatureFlags } from '@/contexts/feature-flag-context';
import { useToast } from '@/hooks/use-toast';

interface SystemBanner {
  id: string;
  title: string;
  message: string;
  type: 'info' | 'warning' | 'error' | 'success';
  active: boolean;
  dismissible: boolean;
  created_at: string;
  expires_at?: string;
}

const bannerTypes = [
  { value: 'info', label: 'Info', icon: Info, color: 'text-blue-600' },
  { value: 'warning', label: 'Warning', icon: AlertTriangle, color: 'text-yellow-600' },
  { value: 'error', label: 'Error', icon: AlertCircle, color: 'text-red-600' },
  { value: 'success', label: 'Success', icon: CheckCircle, color: 'text-green-600' },
];

export default function SystemBanners() {
  const [bannersList, setBannersList] = useState<SystemBanner[]>([]);
  const [loading, setLoading] = useState(true);
  const [newBanner, setNewBanner] = useState<{
    title: string;
    message: string;
    type: 'info' | 'warning' | 'error' | 'success';
    active: boolean;
    dismissible: boolean;
    expires_at: string;
  }>({
    title: '',
    message: '',
    type: 'info',
    active: true,
    dismissible: true,
    expires_at: '',
  });
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  
  const { refreshBanners } = useFeatureFlags();
  const { toast } = useToast();

  useEffect(() => {
    loadBanners();
  }, []);

  const loadBanners = async () => {
    setLoading(true);
    try {
      const response = await fetch('/api/admin/banners');
      if (response.ok) {
        const banners = await response.json();
        setBannersList(banners);
      } else {
        // Mock data for development
        setBannersList([
          {
            id: '1',
            title: 'Scheduled Maintenance',
            message: 'System will undergo maintenance on Jan 25, 2024 from 2:00 AM to 4:00 AM UTC.',
            type: 'warning',
            active: true,
            dismissible: false,
            created_at: '2024-01-20T10:00:00Z',
            expires_at: '2024-01-26T00:00:00Z',
          },
          {
            id: '2',
            title: 'New Features Available',
            message: 'Check out the new portfolio analytics dashboard with advanced charting capabilities.',
            type: 'success',
            active: true,
            dismissible: true,
            created_at: '2024-01-18T14:30:00Z',
          },
          {
            id: '3',
            title: 'API Rate Limits Updated',
            message: 'API rate limits have been increased to improve user experience.',
            type: 'info',
            active: false,
            dismissible: true,
            created_at: '2024-01-15T09:15:00Z',
          },
        ]);
      }
    } catch (error) {
      console.error('Failed to load banners:', error);
      toast({
        title: 'Error loading banners',
        description: 'Failed to fetch system banners from server',
        variant: 'destructive',
      });
    } finally {
      setLoading(false);
    }
  };

  const toggleBanner = async (bannerId: string, active: boolean) => {
    try {
      const response = await fetch(`/api/admin/banners/${bannerId}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ active }),
      });

      if (response.ok) {
        setBannersList(prev => 
          prev.map(banner => 
            banner.id === bannerId 
              ? { ...banner, active }
              : banner
          )
        );
        await refreshBanners();
        toast({
          title: 'Banner updated',
          description: `Banner ${active ? 'activated' : 'deactivated'} successfully`,
        });
      } else {
        throw new Error('Failed to update banner');
      }
    } catch (error) {
      console.error('Failed to toggle banner:', error);
      toast({
        title: 'Error updating banner',
        description: 'Failed to update banner status',
        variant: 'destructive',
      });
    }
  };

  const createBanner = async () => {
    try {
      const response = await fetch('/api/admin/banners', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          ...newBanner,
          expires_at: newBanner.expires_at || null,
        }),
      });

      if (response.ok) {
        const createdBanner = await response.json();
        setBannersList(prev => [...prev, createdBanner]);
        setNewBanner({
          title: '',
          message: '',
          type: 'info',
          active: true,
          dismissible: true,
          expires_at: '',
        });
        setIsDialogOpen(false);
        await refreshBanners();
        toast({
          title: 'Banner created',
          description: 'New system banner created successfully',
        });
      } else {
        throw new Error('Failed to create banner');
      }
    } catch (error) {
      console.error('Failed to create banner:', error);
      toast({
        title: 'Error creating banner',
        description: 'Failed to create new banner',
        variant: 'destructive',
      });
    }
  };

  const deleteBanner = async (bannerId: string) => {
    if (!confirm('Are you sure you want to delete this banner?')) return;

    try {
      const response = await fetch(`/api/admin/banners/${bannerId}`, {
        method: 'DELETE',
      });

      if (response.ok) {
        setBannersList(prev => prev.filter(banner => banner.id !== bannerId));
        await refreshBanners();
        toast({
          title: 'Banner deleted',
          description: 'Banner deleted successfully',
        });
      } else {
        throw new Error('Failed to delete banner');
      }
    } catch (error) {
      console.error('Failed to delete banner:', error);
      toast({
        title: 'Error deleting banner',
        description: 'Failed to delete banner',
        variant: 'destructive',
      });
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-background to-muted/20 p-8">
      <div className="max-w-7xl mx-auto space-y-8">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-12 h-12 bg-primary/10 rounded-lg flex items-center justify-center">
              <Megaphone className="h-6 w-6 text-primary" />
            </div>
            <div>
              <h1 className="text-3xl font-bold">System Banners</h1>
              <p className="text-muted-foreground">Manage application-wide announcements</p>
            </div>
          </div>
          <div className="flex gap-2">
            <Button variant="outline" onClick={loadBanners} disabled={loading}>
              <RefreshCw className={`h-4 w-4 mr-2 ${loading ? 'animate-spin' : ''}`} />
              Refresh
            </Button>
            <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
              <DialogTrigger asChild>
                <Button>
                  <Plus className="h-4 w-4 mr-2" />
                  Create Banner
                </Button>
              </DialogTrigger>
              <DialogContent className="max-w-2xl">
                <DialogHeader>
                  <DialogTitle>Create System Banner</DialogTitle>
                </DialogHeader>
                <div className="space-y-4">
                  <div>
                    <Label htmlFor="title">Title</Label>
                    <Input
                      id="title"
                      value={newBanner.title}
                      onChange={(e) => setNewBanner(prev => ({ ...prev, title: e.target.value }))}
                      placeholder="Scheduled Maintenance"
                    />
                  </div>
                  <div>
                    <Label htmlFor="message">Message</Label>
                    <Textarea
                      id="message"
                      value={newBanner.message}
                      onChange={(e) => setNewBanner(prev => ({ ...prev, message: e.target.value }))}
                      placeholder="Detailed message about the announcement..."
                      rows={3}
                    />
                  </div>
                  <div>
                    <Label htmlFor="type">Type</Label>
                    <Select 
                      value={newBanner.type} 
                      onValueChange={(type: 'info' | 'warning' | 'error' | 'success') => 
                        setNewBanner(prev => ({ ...prev, type }))
                      }
                    >
                      <SelectTrigger>
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        {bannerTypes.map((type) => {
                          const Icon = type.icon;
                          return (
                            <SelectItem key={type.value} value={type.value}>
                              <div className="flex items-center gap-2">
                                <Icon className={`h-4 w-4 ${type.color}`} />
                                {type.label}
                              </div>
                            </SelectItem>
                          );
                        })}
                      </SelectContent>
                    </Select>
                  </div>
                  <div>
                    <Label htmlFor="expires_at">Expires At (Optional)</Label>
                    <Input
                      id="expires_at"
                      type="datetime-local"
                      value={newBanner.expires_at}
                      onChange={(e) => setNewBanner(prev => ({ ...prev, expires_at: e.target.value }))}
                    />
                  </div>
                  <div className="flex items-center space-x-4">
                    <div className="flex items-center space-x-2">
                      <Switch
                        id="active"
                        checked={newBanner.active}
                        onCheckedChange={(active) => setNewBanner(prev => ({ ...prev, active }))}
                      />
                      <Label htmlFor="active">Active</Label>
                    </div>
                    <div className="flex items-center space-x-2">
                      <Switch
                        id="dismissible"
                        checked={newBanner.dismissible}
                        onCheckedChange={(dismissible) => setNewBanner(prev => ({ ...prev, dismissible }))}
                      />
                      <Label htmlFor="dismissible">Dismissible</Label>
                    </div>
                  </div>
                  <Button onClick={createBanner} className="w-full">
                    Create Banner
                  </Button>
                </div>
              </DialogContent>
            </Dialog>
          </div>
        </div>

        {/* Banners Table */}
        <Card>
          <CardHeader>
            <CardTitle>System Banners</CardTitle>
          </CardHeader>
          <CardContent>
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Title</TableHead>
                  <TableHead>Type</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Dismissible</TableHead>
                  <TableHead>Created</TableHead>
                  <TableHead>Expires</TableHead>
                  <TableHead>Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {bannersList.map((banner) => {
                  const bannerType = bannerTypes.find(t => t.value === banner.type);
                  const Icon = bannerType?.icon || Info;
                  
                  return (
                    <TableRow key={banner.id}>
                      <TableCell>
                        <div>
                          <div className="font-medium">{banner.title}</div>
                          <div className="text-sm text-muted-foreground line-clamp-2">
                            {banner.message}
                          </div>
                        </div>
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <Icon className={`h-4 w-4 ${bannerType?.color}`} />
                          <span className="capitalize">{banner.type}</span>
                        </div>
                      </TableCell>
                      <TableCell>
                        <Badge variant={banner.active ? 'default' : 'secondary'}>
                          {banner.active ? 'Active' : 'Inactive'}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        <Badge variant={banner.dismissible ? 'outline' : 'secondary'}>
                          {banner.dismissible ? 'Yes' : 'No'}
                        </Badge>
                      </TableCell>
                      <TableCell>
                        {new Date(banner.created_at).toLocaleDateString()}
                      </TableCell>
                      <TableCell>
                        {banner.expires_at 
                          ? new Date(banner.expires_at).toLocaleDateString()
                          : 'Never'
                        }
                      </TableCell>
                      <TableCell>
                        <div className="flex items-center gap-2">
                          <Switch
                            checked={banner.active}
                            onCheckedChange={(active) => toggleBanner(banner.id, active)}
                          />
                          <Button variant="ghost" size="sm">
                            <Edit className="h-4 w-4" />
                          </Button>
                          <Button 
                            variant="ghost" 
                            size="sm"
                            onClick={() => deleteBanner(banner.id)}
                          >
                            <Trash2 className="h-4 w-4 text-destructive" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  );
                })}
              </TableBody>
            </Table>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}