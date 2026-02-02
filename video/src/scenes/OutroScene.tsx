import React from 'react';
import {AbsoluteFill, useCurrentFrame, interpolate, Easing} from 'remotion';
import {COLORS, FONTS, SIZES} from '../constants';

interface OutroSceneProps {
  startFrame: number;
}

export const OutroScene: React.FC<OutroSceneProps> = ({startFrame}) => {
  const frame = useCurrentFrame();
  const relativeFrame = frame - startFrame;
  
  const sceneOpacity = interpolate(relativeFrame, [0, 30], [0, 1], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
    easing: Easing.out(Easing.ease),
  });
  
  const logoScale = interpolate(relativeFrame, [0, 40], [0.8, 1], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
    easing: Easing.out(Easing.back(1.7)),
  });
  
  const ctaOpacity = interpolate(relativeFrame, [40, 70], [0, 1], {
    extrapolateLeft: 'clamp',
    extrapolateRight: 'clamp',
    easing: Easing.out(Easing.ease),
  });
  
  const installOpacity = interpolate(relativeFrame, [60, 90], [0, 1], {
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
      opacity: sceneOpacity,
    }}>
      {/* Logo */}
      <div style={{
        display: 'flex',
        alignItems: 'center',
        gap: '24px',
        transform: `scale(${logoScale})`,
        marginBottom: '40px',
      }}>
        <div style={{
          width: '100px',
          height: '100px',
          backgroundColor: COLORS.accentPrimary,
          borderRadius: '20px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          boxShadow: '0 8px 30px rgba(167, 206, 96, 0.4)',
        }}>
          <span style={{
            fontSize: '60px',
            fontWeight: 'bold',
            color: COLORS.text,
          }}>
            M
          </span>
        </div>
        
        <div>
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
            marginLeft: '12px',
          }}>
            cli
          </span>
        </div>
      </div>
      
      {/* CTA */}
      <div style={{
        opacity: ctaOpacity,
        textAlign: 'center',
        marginBottom: '40px',
      }}>
        <p style={{
          fontSize: SIZES.body,
          color: COLORS.textLight,
          margin: '0 0 24px 0',
        }}>
          The powerful CLI for Mochi.cards
        </p>
        
        <div style={{
          display: 'inline-flex',
          alignItems: 'center',
          gap: '12px',
          padding: '16px 32px',
          backgroundColor: COLORS.accentPrimary,
          borderRadius: '100px',
          boxShadow: '0 4px 20px rgba(167, 206, 96, 0.3)',
        }}>
          <span style={{
            fontSize: SIZES.body,
            color: COLORS.text,
            fontWeight: '600',
          }}>
            Get Started Today
          </span>
        </div>
      </div>
      
      {/* Install Command */}
      <div style={{
        opacity: installOpacity,
      }}>
        <div style={{
          backgroundColor: COLORS.terminalBg,
          padding: '20px 32px',
          borderRadius: '12px',
          fontFamily: FONTS.mono,
          fontSize: SIZES.code,
          color: COLORS.terminalText,
          boxShadow: '0 4px 20px rgba(0, 0, 0, 0.1)',
        }}>
          <span style={{ color: COLORS.accentPrimary }}>$</span>
          {' '}
          <span>curl -fsSL</span>
          {' '}
          <span style={{ color: COLORS.terminalGreen }}>github.com/nerveband/mochi-cli</span>
          {' '}
          <span>| bash</span>
        </div>
      </div>
      
      {/* Footer */}
      <div style={{
        position: 'absolute',
        bottom: '40px',
        opacity: interpolate(relativeFrame, [100, 130], [0, 1], {
          extrapolateLeft: 'clamp',
          extrapolateRight: 'clamp',
        }),
      }}>
        <p style={{
          fontSize: SIZES.small,
          color: COLORS.textMuted,
          margin: 0,
        }}>
          github.com/nerveband/mochi-cli
        </p>
      </div>
    </AbsoluteFill>
  );
};
