import React from 'react';
import {AbsoluteFill, useCurrentFrame, interpolate, Easing} from 'remotion';
import {COLORS, FONTS, SIZES} from '../constants';

export const IntroScene: React.FC = () => {
  const frame = useCurrentFrame();
  
  // Animation timing
  const titleOpacity = interpolate(frame, [0, 30], [0, 1], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
    easing: Easing.out(Easing.ease),
  });
  
  const titleY = interpolate(frame, [0, 30], [50, 0], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
    easing: Easing.out(Easing.ease),
  });
  
  const subtitleOpacity = interpolate(frame, [20, 50], [0, 1], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
    easing: Easing.out(Easing.ease),
  });
  
  const taglineOpacity = interpolate(frame, [40, 70], [0, 1], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
    easing: Easing.out(Easing.ease),
  });
  
  return (
    <AbsoluteFill style={{
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      alignItems: 'center',
      backgroundColor: COLORS.background,
    }}>
      {/* Logo/Title */}
      <div style={{
        display: 'flex',
        alignItems: 'center',
        gap: '24px',
        opacity: titleOpacity,
        transform: `translateY(${titleY}px)`,
      }}>
        {/* Mochi icon representation */}
        <div style={{
          width: '80px',
          height: '80px',
          backgroundColor: COLORS.accentPrimary,
          borderRadius: '16px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          boxShadow: '0 4px 20px rgba(167, 206, 96, 0.3)',
        }}>
          <span style={{
            fontSize: '48px',
            fontWeight: 'bold',
            color: COLORS.text,
          }}>
            M
          </span>
        </div>
        
        <div style={{
          display: 'flex',
          alignItems: 'baseline',
          gap: '12px',
        }}>
          <span style={{
            fontSize: SIZES.title,
            fontWeight: '700',
            color: COLORS.text,
            letterSpacing: '-1px',
          }}>
            mochi
          </span>
          <span style={{
            fontSize: SIZES.subtitle,
            fontWeight: '300',
            color: COLORS.textLight,
          }}>
            cli
          </span>
        </div>
      </div>
      
      {/* Subtitle */}
      <div style={{
        marginTop: '32px',
        opacity: subtitleOpacity,
      }}>
        <p style={{
          fontSize: SIZES.body,
          color: COLORS.textLight,
          textAlign: 'center',
          maxWidth: '600px',
          lineHeight: '1.6',
        }}>
          Command-line interface for Mochi.cards
        </p>
      </div>
      
      {/* Tagline */}
      <div style={{
        marginTop: '24px',
        opacity: taglineOpacity,
        display: 'flex',
        alignItems: 'center',
        gap: '8px',
        padding: '12px 24px',
        backgroundColor: COLORS.accentSecondary,
        borderRadius: '100px',
      }}>
        <span style={{
          fontSize: SIZES.small,
          color: COLORS.text,
          fontWeight: '500',
        }}>
          Built for LLMs & Automation
        </span>
      </div>
      
      {/* Decorative elements */}
      <div style={{
        position: 'absolute',
        bottom: '60px',
        left: '50%',
        transform: 'translateX(-50%)',
        display: 'flex',
        gap: '8px',
      }}>
        {[0, 1, 2].map((i) => (
          <div
            key={i}
            style={{
              width: '8px',
              height: '8px',
              borderRadius: '50%',
              backgroundColor: COLORS.accentPrimary,
              opacity: interpolate(frame, [60 + i * 10, 80 + i * 10], [0, 1], {
                extrapolateLeft: 'clamp',
                extrapolateRight: 'clamp',
              }),
            }}
          />
        ))}
      </div>
    </AbsoluteFill>
  );
};
