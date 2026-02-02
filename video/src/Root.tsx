import React from 'react';
import {Composition} from 'remotion';
import {MochiCLIDemo} from './scenes/MochiCLIDemo';
import {MochiCliPromo} from './MochiCliPromo';
import {COMPOSITION_CONFIG} from '../remotion.config';

export const Root: React.FC = () => {
  return (
    <>
      <Composition
        id="MochiCLIDemo"
        component={MochiCLIDemo}
        durationInFrames={COMPOSITION_CONFIG.durationInFrames}
        fps={COMPOSITION_CONFIG.fps}
        width={COMPOSITION_CONFIG.width}
        height={COMPOSITION_CONFIG.height}
      />
      <Composition
        id="MochiCliPromo"
        component={MochiCliPromo}
        durationInFrames={450} // 15 seconds at 30fps
        fps={30}
        width={1920}
        height={1080}
      />
    </>
  );
};
