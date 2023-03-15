import React, {useEffect, useRef} from 'react';
import styles from '../../pages/Call/Call.module.scss';
import {MutedAudio, MutedVideo, SignalIcon, VideoPlaceholder} from '../../assets';
import {Box} from '@chakra-ui/react';
import {useBreakpoints} from '../../hooks/useBreakpoints';

const Video = ({participant, stream, isCurrentUser, mediaState, message, subtitlesEnabled}) => {
  const container = useRef();
  const videoElement = useRef();
  const {isMobile} = useBreakpoints();

  useEffect(() => {
    if (videoElement.current && stream && !videoElement.current.srcObject) {
      videoElement.current.srcObject = stream;
      // viewer mode safari not working autoplay without get user media
      videoElement.current.play();
    }
  });

  return (
    <Box
      p={isMobile ? '4px' : '8px'}
      height={isMobile ? 'calc(100% - 8px)' : 'calc(100% - 16px)'}
    >
      <div
        ref={container}
        className={styles.streamContainer}
      >
        <video
          ref={videoElement}
          id={'stream_' + participant.streamID}
          autoPlay
          muted={isCurrentUser}
          playsInline
          style={{
            opacity: mediaState?.video ? '1' : '0',
            transform: isCurrentUser ? 'scale(-1, 1)' : 'none',
          }}
        />

        <div className={styles.badge}>
          <img
            src={SignalIcon}
            alt={'signal strength icon'}
          />
          <p id={'badgeText-' + participant.streamID}>
            {participant.name}{participant.isHost ? ' (Host)' : ''}
          </p>
        </div>


        <div className={styles.controlsBox}>
          {!mediaState?.audio && (
            <img
              src={MutedAudio}
              alt={'muted audio icon'}
            />
          )}

          {!mediaState?.video && (
            <img
              src={MutedVideo}
              alt={'muted video icon'}
            />
          )}
        </div>

        <img
          src={VideoPlaceholder}
          alt={'Video placeholder'}
          className={styles.streamPlaceholder}
        />

        {subtitlesEnabled && (
          <>
            <Box
              opacity={message ? 1 : 0}
              className={styles.messageOverlay}
            />

            {message && (
              <div className={styles.message}>
                <p>
                  {message}
                </p>
              </div>
            )}
          </>
        )}

      </div>
    </Box>
  );
};

export default Video;
