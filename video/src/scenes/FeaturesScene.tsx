import React from 'react';
import {AbsoluteFill, useCurrentFrame, interpolate, Easing} from 'remotion';
import {COLORS, FONTS, SIZES} from '../constants';

interface FeaturesSceneProps {
  startFrame: number;
}

interface Feature {
  icon: string;
  title: string;
  description: string;
}

const features: Feature[] = [
  {
    icon: '🤖',
    title: 'LLM-Optimized',
    description: 'JSON output, quiet mode, field extraction',
  },
  {
    icon: '⚡',
    title: 'Multiple Formats',
    description: 'JSON, Table, Markdown, Compact',
  },
  {
    icon: '🔑',
    title: 'Multi-Profile',
    description: 'Switch between API keys easily',
  },
  {
    icon: '🔄',
    title: 'Self-Updating',
    description: 'Always stay current with latest features',
  },
  {
    icon: '💻',
    title: 'Pipe-Friendly',
    description: 'Read from stdin, perfect for scripts',
  },
  {
    icon: '🔒',
    title: 'Secure',
    description: 'Environment variables or secure config',
  },
];

export const FeaturesScene: React.FC<FeaturesSceneProps> = ({startFrame}) => {
  const frame = useCurrentFrame();
  const relativeFrame = frame - startFrame;
  
  const sceneOpacity = interpolate(relativeFrame, [0, 30], [0, 1], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
    easing: Easing.out(Easing.ease),
  });
  
  const sceneExitOpacity = interpolate(relativeFrame, [270, 300], [1, 0], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
    easing: Easing.in(Easing.ease),
  });
  
  return (
    <AbsoluteFill style={{
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      alignItems: 'center',
      backgroundColor: COLORS.background,
      opacity: Math.min(sceneOpacity, sceneExitOpacity),
      padding: '80px',
    }}>
      {/* Title */}
      <div style={{
        marginBottom: '60px',
        opacity: interpolate(relativeFrame, [0, 30], [0, 1], {
          extrapolateLeft: 'clamp',
          extrapolateRight: 'clamp',
        }),
      }}>
        <h2 style={{
          fontSize: SIZES.subtitle,
          color: COLORS.text,
          fontWeight: '600',
          margin: 0,
          textAlign: 'center',
        }}>
          Built for Modern Workflows
        </h2>
      </div>
      
      {/* Features Grid */}
      <div style={{
        display: 'grid',
        gridTemplateColumns: 'repeat(3, 1fr)',
        gap: '40px',
        maxWidth: '1200px',
      }}>
        {features.map((feature, idx) => {
          const featureDelay = idx * 15;
          const featureOpacity = interpolate(
            relativeFrame,
            [featureDelay, featureDelay + 30],
            [0, 1],
            {
              extrapolateLeft: 'clamp',
              extrapolateRight: 'clamp',
              easing: Easing.out(Easing.ease),
            }
          );
          
          const featureY = interpolate(
            relativeFrame,
            [featureDelay, featureDelay + 30],
            [30, 0],
            {
              extrapolateLeft: 'clamp',
              extrapolateRight: 'clamp',
              easing: Easing.out(Easing.ease),
            }
          );
          
          return (
            <div
              key={idx}
              style={{
                opacity: featureOpacity,
                transform: `translateY(${featureY}px)`,
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                textAlign: 'center',
                padding: '30px',
                backgroundColor: COLORS.cardBg,
                borderRadius: '16px',
                boxShadow: '0 4px 20px rgba(0, 0, 0, 0.05)',
              }}
            >
              <div style={{
                fontSize: '48px',
                marginBottom: '16px',
              }}>
                {feature.icon}
              </div>
              <h3 style={{
                fontSize: SIZES.body,
                color: COLORS.text,
                fontWeight: '600',
                margin: '0 0 8px 0',
              }}>
                {feature.title}
              </h3>
              <p style={{
                fontSize: SIZES.small,
                color: COLORS.textLight,
                margin: 0,
                lineHeight: '1.5',
              }}>
                {feature.description}
              </p>
            </div>
          );
        })}
      </div>
    </AbsoluteFill>
  );
};
