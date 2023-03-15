import React from 'react';
import styles from './VideoControls.module.scss';
import SourceControl from '../SourceControl/SourceControl';
import {HangUpIcon} from '../../assets';
import ParticipantsBadge from '../ParticipantsBadge/ParticipantsBadge';
import classNames from 'classnames';
import SubtitlesControl from '../SourceControl/SubtitlesControl';

const VideoControls = ({
  devices,
  onHangUp,
  videoEnabled,
  audioEnabled,
  toggleAudio,
  toggleVideo,
  onDeviceChange,
  selectedVideoId,
  selectedAudioId,
  isCall,
  participantsCount,
  noPublish,
  subtitlesEnabled,
  toggleSubtitles
}) => {
  return (
    <div className={classNames(styles.container, isCall && styles.containerCall)}>
      {isCall && (
        <div className={styles.participants}>
          <ParticipantsBadge count={participantsCount}/>
        </div>
      )}

      {!noPublish && (
        <>
          <SourceControl
            onChange={(deviceID) => onDeviceChange('audio', deviceID)}
            devices={devices.filter(d => d.kind === 'audioinput')}
            enabled={audioEnabled}
            toggleMute={toggleAudio}
            selected={selectedAudioId}
            isCall={isCall}
          />

          <SourceControl
            onChange={(deviceID) => onDeviceChange('video', deviceID)}
            devices={devices.filter(d => d.kind === 'videoinput')}
            isVideo
            enabled={videoEnabled}
            toggleMute={toggleVideo}
            selected={selectedVideoId}
            isCall={isCall}
          />

          {isCall && (
            <SubtitlesControl
              enabled={subtitlesEnabled}
              toggleMute={toggleSubtitles}
              isCall={isCall}
            />
          )}
        </>
      )}

      {onHangUp && (
        <button
          className={styles.hangup}
          onClick={onHangUp}
        >
          <img
            src={HangUpIcon}
            alt={'hangup icon'}
          />
          <span>End Meeting</span>
        </button>
      )}
    </div>
  );
};

export default VideoControls;
