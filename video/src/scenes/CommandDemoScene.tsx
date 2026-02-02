import React from 'react';
import {AbsoluteFill, useCurrentFrame, interpolate, Easing} from 'remotion';
import {COLORS, FONTS, SIZES} from '../constants';

interface CommandDemoSceneProps {
  startFrame: number;
}

interface TerminalLine {
  text: string;
  color?: string;
  delay: number;
  typingSpeed?: number;
}

export const CommandDemoScene: React.FC<CommandDemoSceneProps> = ({startFrame}) => {
  const frame = useCurrentFrame();
  const relativeFrame = frame - startFrame;
  
  // Scene opacity for transition
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
  
  // Commands to demonstrate
  const commands: TerminalLine[] = [
    { text: '$ mochi deck list --format table', delay: 30, color: COLORS.terminalText },
    { text: 'ID          NAME           SORT  ARCHIVED', delay: 70, color: COLORS.textMuted },
    { text: 'ABCD1234    Spanish        1', delay: 80, color: COLORS.terminalText },
    { text: 'EFGH5678    Programming    2', delay: 90, color: COLORS.terminalText },
    { text: '', delay: 110 },
    { text: '$ mochi card create --deck ABCD1234 --content "# Hola"', delay: 130, color: COLORS.terminalText },
    { text: 'Card created: XYZ999', delay: 190, color: COLORS.accentPrimary },
    { text: '', delay: 210 },
    { text: '$ mochi due list', delay: 230, color: COLORS.terminalText },
    { text: '3 cards due today', delay: 270, color: COLORS.accentSecondary },
  ];
  
  // Calculate visible text for each line
  const getVisibleText = (line: TerminalLine, idx: number) => {
    const lineFrame = relativeFrame - line.delay;
    if (lineFrame < 0) return '';
    
    const text = line.text;
    const duration = text.length * 2; // 2 frames per character
    
    if (lineFrame > duration) return text;
    
    return text.slice(0, Math.floor(lineFrame / 2));
  };
  
  const getCursorOpacity = (line: TerminalLine, idx: number) => {
    const lineFrame = relativeFrame - line.delay;
    const text = line.text;
    const duration = text.length * 2;
    
    if (lineFrame < 0) return 0;
    if (lineFrame > duration + 20) return 0;
    
    // Blinking cursor
    return interpolate(lineFrame - duration, [0, 5, 10, 15], [0, 1, 1, 0], {
      extrapolateLeft: 'clamp',
      extrapolateRight: 'clamp',
    });
  };
  
  const visibleLines = commands.filter((line, idx) => {
    return relativeFrame >= line.delay;
  });
  
  return (
    <AbsoluteFill style={{
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      alignItems: 'center',
      backgroundColor: COLORS.background,
      opacity: Math.min(sceneOpacity, sceneExitOpacity),
    }}>
      {/* Section Title */}
      <div style={{
        position: 'absolute',
        top: '80px',
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
        }}>
          Powerful CLI Commands
        </h2>
      </div>
      
      {/* Terminal Window */}
      <div style={{
        width: '900px',
        backgroundColor: COLORS.terminalBg,
        borderRadius: '12px',
        overflow: 'hidden',
        boxShadow: '0 20px 60px rgba(0, 0, 0, 0.15)',
        border: `1px solid ${COLORS.border}`,
      }}>
        {/* Terminal Header */}
        <div style={{
          height: '40px',
          backgroundColor: '#2a2a2a',
          display: 'flex',
          alignItems: 'center',
          padding: '0 16px',
          gap: '8px',
        }}>
          <div style={{ width: '12px', height: '12px', borderRadius: '50%', backgroundColor: '#ff5f57' }} />
          <div style={{ width: '12px', height: '12px', borderRadius: '50%', backgroundColor: '#febc2e' }} />
          <div style={{ width: '12px', height: '12px', borderRadius: '50%', backgroundColor: '#28c840' }} />
          <span style={{
            marginLeft: 'auto',
            fontSize: '14px',
            color: COLORS.textMuted,
            fontFamily: FONTS.mono,
          }}>
            mochi-cli
          </span>
        </div>
        
        {/* Terminal Content */}
        <div style={{
          padding: '24px',
          fontFamily: FONTS.mono,
          fontSize: SIZES.code,
          lineHeight: '1.8',
          minHeight: '350px',
        }}>
          {commands.map((line, idx) => {
            const visibleText = getVisibleText(line, idx);
            const cursorOpacity = getCursorOpacity(line, idx);
            const lineOpacity = relativeFrame >= line.delay ? 1 : 0;
            
            return (
              <div
                key={idx}
                style={{
                  opacity: lineOpacity,
                  color: line.color || COLORS.terminalText,
                  whiteSpace: 'pre-wrap',
                  wordBreak: 'break-all',
                }}
              >
                {visibleText}
                <span style={{
                  display: 'inline-block',
                  width: '10px',
                  height: '20px',
                  backgroundColor: COLORS.accentPrimary,
                  opacity: cursorOpacity,
                  marginLeft: '2px',
                  verticalAlign: 'middle',
                }} />
              </div>
            );
          })}
        </div>
      </div>
      
      {/* Bottom hint */}
      <div style={{
        position: 'absolute',
        bottom: '60px',
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
          Manage decks, cards, and templates from the command line
        </p>
      </div>
    </AbsoluteFill>
  );
};
