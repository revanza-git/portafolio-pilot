import { useEffect } from 'react';
import { Navigate } from 'react-router-dom';
import { TrendingUp, Shield, Zap, BarChart3, ArrowRight } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { WalletConnectButton } from '@/components/wallet/wallet-connect-button';
import { useWalletStore } from '@/stores/wallet';

const Index = () => {
  const { isConnected } = useWalletStore();

  // Redirect to dashboard if already connected
  if (isConnected) {
    return <Navigate to="/dashboard" replace />;
  }

  const features = [
    {
      icon: BarChart3,
      title: "Portfolio Tracking",
      description: "Monitor your DeFi positions and track portfolio performance across multiple protocols"
    },
    {
      icon: Shield,
      title: "Security Management", 
      description: "Review and revoke token approvals to protect your assets from unnecessary permissions"
    },
    {
      icon: Zap,
      title: "Token Swaps",
      description: "Execute trades with the best rates using integrated DEX aggregators"
    }
  ];

  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-muted/20">
      {/* Hero Section */}
      <div className="relative overflow-hidden">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="relative z-10 pt-20 pb-16 md:pt-28 md:pb-20">
            <div className="text-center">
              {/* Logo */}
              <div className="flex justify-center mb-8">
                <div className="w-16 h-16 bg-gradient-primary rounded-2xl flex items-center justify-center shadow-glow">
                  <TrendingUp className="h-8 w-8 text-primary-foreground" />
                </div>
              </div>

              {/* Title */}
              <h1 className="text-5xl md:text-7xl font-bold tracking-tight mb-8">
                <span className="bg-gradient-primary bg-clip-text text-transparent">
                  DeFi Portfolio
                </span>
                <br />
                <span className="text-foreground">Management</span>
              </h1>

              {/* Subtitle */}
              <p className="text-xl md:text-2xl text-muted-foreground max-w-3xl mx-auto mb-12 leading-relaxed">
                Take control of your decentralized finance investments with comprehensive 
                portfolio tracking, security management, and seamless token swapping.
              </p>

              {/* CTA */}
              <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-16">
                <WalletConnectButton />
                <Button variant="outline" size="lg" className="group">
                  Learn More
                  <ArrowRight className="ml-2 h-4 w-4 group-hover:translate-x-1 transition-transform" />
                </Button>
              </div>

              {/* Stats */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-4xl mx-auto">
                <div className="text-center">
                  <div className="text-3xl font-bold text-primary mb-2">$2.1B+</div>
                  <div className="text-muted-foreground">Total Value Tracked</div>
                </div>
                <div className="text-center">
                  <div className="text-3xl font-bold text-primary mb-2">50K+</div>
                  <div className="text-muted-foreground">Active Users</div>
                </div>
                <div className="text-center">
                  <div className="text-3xl font-bold text-primary mb-2">99.9%</div>
                  <div className="text-muted-foreground">Uptime</div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Features Section */}
      <div className="py-20 bg-gradient-card">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold mb-4">
              Everything you need for DeFi management
            </h2>
            <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
              Professional tools designed for both beginners and experienced DeFi users
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {features.map((feature, index) => (
              <Card key={index} className="border-0 shadow-card bg-gradient-card">
                <CardContent className="p-8 text-center">
                  <div className="w-12 h-12 bg-gradient-primary rounded-xl flex items-center justify-center mx-auto mb-6 shadow-soft">
                    <feature.icon className="h-6 w-6 text-primary-foreground" />
                  </div>
                  <h3 className="text-xl font-semibold mb-4">{feature.title}</h3>
                  <p className="text-muted-foreground leading-relaxed">
                    {feature.description}
                  </p>
                </CardContent>
              </Card>
            ))}
          </div>
        </div>
      </div>

      {/* Footer */}
      <div className="py-12 text-center text-muted-foreground">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <p>&copy; 2024 DeFi Portfolio. Built with modern web technologies.</p>
        </div>
      </div>
    </div>
  );
};

export default Index;
