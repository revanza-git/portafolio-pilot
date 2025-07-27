import { useState } from 'react';
import { format } from 'date-fns';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { AlertCircle, CheckCircle, Clock, RefreshCw } from 'lucide-react';
import { useAlertsStore } from '@/stores/alerts';
import { AlertHistory } from '@/stores/alerts';

export function AlertHistoryComponent() {
  const { history } = useAlertsStore();
  const [statusFilter, setStatusFilter] = useState<'all' | 'sent' | 'failed' | 'pending'>('all');
  const [sortBy, setSortBy] = useState<'recent' | 'oldest'>('recent');

  const filteredHistory = history
    .filter(entry => statusFilter === 'all' || entry.status === statusFilter)
    .sort((a, b) => {
      return sortBy === 'recent' 
        ? b.triggeredAt - a.triggeredAt
        : a.triggeredAt - b.triggeredAt;
    });

  const getStatusIcon = (status: AlertHistory['status']) => {
    switch (status) {
      case 'sent':
        return <CheckCircle className="h-4 w-4 text-profit" />;
      case 'failed':
        return <AlertCircle className="h-4 w-4 text-destructive" />;
      case 'pending':
        return <Clock className="h-4 w-4 text-warning" />;
      default:
        return null;
    }
  };

  const getStatusVariant = (status: AlertHistory['status']) => {
    switch (status) {
      case 'sent':
        return 'default';
      case 'failed':
        return 'destructive';
      case 'pending':
        return 'secondary';
      default:
        return 'outline';
    }
  };

  const formatValue = (value: number, token: string) => {
    if (token.includes('APR') || token.includes('Pool')) {
      return `${value.toFixed(2)}%`;
    }
    return `$${value.toLocaleString()}`;
  };

  return (
    <Card className="bg-gradient-card shadow-card border-0">
      <CardHeader>
        <div className="flex justify-between items-center">
          <CardTitle>Alert History</CardTitle>
          <div className="flex items-center space-x-2">
            <Select value={statusFilter} onValueChange={(value: any) => setStatusFilter(value)}>
              <SelectTrigger className="w-32">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All Status</SelectItem>
                <SelectItem value="sent">Sent</SelectItem>
                <SelectItem value="failed">Failed</SelectItem>
                <SelectItem value="pending">Pending</SelectItem>
              </SelectContent>
            </Select>
            <Select value={sortBy} onValueChange={(value: any) => setSortBy(value)}>
              <SelectTrigger className="w-32">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="recent">Recent First</SelectItem>
                <SelectItem value="oldest">Oldest First</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>
      </CardHeader>
      <CardContent>
        {filteredHistory.length === 0 ? (
          <div className="text-center py-8">
            <Clock className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
            <h3 className="text-lg font-medium mb-2">No alert history</h3>
            <p className="text-muted-foreground">
              Alert notifications will appear here once triggered
            </p>
          </div>
        ) : (
          <div className="rounded-md border border-border">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Triggered</TableHead>
                  <TableHead>Alert</TableHead>
                  <TableHead>Condition</TableHead>
                  <TableHead>Value</TableHead>
                  <TableHead>Channel</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Retries</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {filteredHistory.map((entry) => (
                  <TableRow key={entry.id}>
                    <TableCell>
                      <div className="text-sm">
                        {format(new Date(entry.triggeredAt), 'MMM dd, HH:mm')}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        {format(new Date(entry.triggeredAt), 'yyyy')}
                      </div>
                    </TableCell>
                    <TableCell className="font-medium">{entry.token}</TableCell>
                    <TableCell>
                      <div className="flex items-center space-x-2">
                        <span className="text-sm">
                          {entry.condition} {formatValue(entry.threshold, entry.token)}
                        </span>
                      </div>
                    </TableCell>
                    <TableCell>
                      <span className="font-medium">
                        {formatValue(entry.value, entry.token)}
                      </span>
                    </TableCell>
                    <TableCell>
                      <Badge variant="outline" className="capitalize">
                        {entry.channel}
                      </Badge>
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center space-x-2">
                        {getStatusIcon(entry.status)}
                        <Badge variant={getStatusVariant(entry.status)} className="capitalize">
                          {entry.status}
                        </Badge>
                      </div>
                      {entry.error && (
                        <div className="text-xs text-destructive mt-1" title={entry.error}>
                          {entry.error.length > 30 ? `${entry.error.substring(0, 30)}...` : entry.error}
                        </div>
                      )}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center space-x-1">
                        {entry.retryCount > 0 && (
                          <RefreshCw className="h-3 w-3 text-muted-foreground" />
                        )}
                        <span className="text-sm">{entry.retryCount}</span>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        )}
      </CardContent>
    </Card>
  );
}