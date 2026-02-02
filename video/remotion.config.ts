import {Config} from '@remotion/cli/config';

Config.setVideoImageFormat('jpeg');
Config.setOverwriteOutput(true);

export const COMPOSITION_CONFIG = {
  width: 1920,
  height: 1080,
  fps: 60,
  durationInFrames: 900, // 15 seconds
};
