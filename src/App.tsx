import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { ThemeProvider } from "@/components/providers/theme-provider";
import Index from "./pages/Index";
import Dashboard from "./pages/Dashboard";
import Yield from "./pages/Yield";
import Bridge from "./pages/Bridge";
import Analytics from "./pages/Analytics";
import Transactions from "./pages/Transactions";
import Approvals from "./pages/Approvals";
import Swap from "./pages/Swap";
import Alerts from "./pages/Alerts";
import Watchlist from "./pages/Watchlist";
import NotFound from "./pages/NotFound";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <ThemeProvider>
      <TooltipProvider>
        <Toaster />
        <Sonner />
        <BrowserRouter>
          <Routes>
            <Route path="/" element={<Index />} />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/yield" element={<Yield />} />
            <Route path="/bridge" element={<Bridge />} />
            <Route path="/analytics" element={<Analytics />} />
            <Route path="/transactions" element={<Transactions />} />
            <Route path="/approvals" element={<Approvals />} />
            <Route path="/swap" element={<Swap />} />
            <Route path="/alerts" element={<Alerts />} />
            <Route path="/watchlist" element={<Watchlist />} />
            {/* ADD ALL CUSTOM ROUTES ABOVE THE CATCH-ALL "*" ROUTE */}
            <Route path="*" element={<NotFound />} />
          </Routes>
        </BrowserRouter>
      </TooltipProvider>
    </ThemeProvider>
  </QueryClientProvider>
);

export default App;
