import React from 'react';
import {AbsoluteFill, useCurrentFrame, useVideoConfig, interpolate, Easing} from 'remotion';
import {COLORS, FONTS, SIZES} from '../constants';
import {IntroScene} from './IntroScene';
import {CommandDemoScene} from './CommandDemoScene';
import {FeaturesScene} from './FeaturesScene';
import {OutroScene} from './OutroScene';

export const MochiCLIDemo: React.FC = () => {
  const frame = useCurrentFrame();
  const {fps} = useVideoConfig();
  
  // Scene timing (in frames)
  const introEnd = 120; // 2 seconds
  const commandsEnd = 420; // 5 seconds of demos
  const featuresEnd = 720; // 5 seconds
  const outroEnd = 900; // 3 seconds
  
  return (
    <AbsoluteFill style={{
      backgroundColor: COLORS.background,
      fontFamily: FONTS.primary,
      color: COLORS.text,
    }}>
      {/* Scene transitions */}
      {frame < introEnd && <IntroScene />}
      {frame >= introEnd - 30 && frame < commandsEnd && <CommandDemoScene startFrame={introEnd} />}
      {frame >= commandsEnd - 30 && frame < featuresEnd && <FeaturesScene startFrame={commandsEnd} />}
      {frame >= featuresEnd - 30 && <OutroScene startFrame={featuresEnd} />}
    </AbsoluteFill>
  );
};
