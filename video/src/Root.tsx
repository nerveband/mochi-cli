import React from 'react';
import {Composition} from 'remotion';
import {MochiCLIDemo} from './scenes/MochiCLIDemo';
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
    </>
  );
};
