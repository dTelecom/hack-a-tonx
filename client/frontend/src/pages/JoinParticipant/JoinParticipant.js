import React, {useCallback, useEffect, useRef, useState} from 'react';
import {Header} from '../../components/Header/Header';
import styles from './JoinParticipant.module.scss';
import {observer} from 'mobx-react';
import {useLocation, useNavigate, useParams} from 'react-router-dom';
import {Container} from '../../components/Container/Container';
import {Box, Flex} from '@chakra-ui/react';
import Input from '../../components/Input/Input';
import VideoControls from '../../components/VideoControls/VideoControls';
import Footer from '../../components/Footer/Footer';
import {createVideoElement, hideMutedBadge, showMutedBadge} from '../Call/utils';
import {useMediaConstraints} from '../../hooks/useMediaConstraints';
import ParticipantsBadge from '../../components/ParticipantsBadge/ParticipantsBadge';
import {useBreakpoints} from '../../hooks/useBreakpoints';
import {loadDevices} from '../../utils/loadDevices';
import {FaceIcon} from '../../assets';
import {Button} from '../../components/Button/Button';

const JoinParticipant = () => {
  const navigate = useNavigate();
  const {isMobile} = useBreakpoints();
  const [name, setName] = useState('');
  const [hasVideo, setHasVideo] = useState(false);
  const [devices, setDevices] = useState([]);
  const location = useLocation();
  const [room] = useState(location.state?.room);
  const {
    constraints,
    onDeviceChange,
    onMediaToggle,
    audioEnabled,
    videoEnabled,
    selectedAudioId,
    selectedVideoId,
    constraintsState,
  } = useMediaConstraints();
  const {sid} = useParams();
  const videoContainer = useRef();
  const localVideo = useRef();

  useEffect(() => {
    void loadMedia(constraints);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const loadMedia = useCallback(async (config) => {
    console.log('[loadMedia]', config);
    try {
      const stream = await navigator.mediaDevices.getUserMedia(config);
      void loadDevices(setDevices);

      if (!selectedVideoId && !selectedAudioId) {
        // set initial devices
        stream.getTracks().forEach(track => {
            const deviceId = track.getSettings().deviceId;
            onDeviceChange(track.kind, deviceId);
          }
        );
      }

      if (!videoContainer.current) {
        setTimeout(() => loadMedia(config), 200);
      } else {
        localVideo.current = stream;
        const video = createVideoElement({
          media: stream,
          muted: true,
          hideBadge: true,
          style: {width: '100%', height: '100%', transform: 'scale(-1, 1)'},
          audio: !!config.audio,
          video: !!config.video,
        });
        video.style.transform = 'scale(-1, 1)';

        videoContainer.current.innerHTML = '';
        videoContainer.current.appendChild(video);
        setHasVideo(true);
      }
    } catch
      (err) {
      console.error(err);
    }
  }, [onDeviceChange, selectedAudioId, selectedVideoId]);

  const onDeviceSelect = useCallback((type, deviceId) => {
    const constraints = onDeviceChange(type, deviceId);
    void loadMedia(constraints);
  }, [loadMedia, onDeviceChange]);

  function toggleAudio() {
    if (localVideo.current) {
      const track = localVideo.current.getAudioTracks()[0];
      if (!track) {
        onDeviceSelect('audio', true);
        return;
      }
      track.enabled = !audioEnabled;
      if (!audioEnabled) {
        hideMutedBadge('audio', localVideo.current.id);
      } else {
        showMutedBadge('audio', localVideo.current.id);
      }
      onMediaToggle('audio');
    }
  }

  function toggleVideo() {
    if (localVideo.current) {
      const prevState = videoEnabled;
      const track = localVideo.current.getVideoTracks()[0];
      if (!track) {
        onDeviceSelect('video', true);
        return;
      }
      track.enabled = !prevState;
      if (!prevState) {
        hideMutedBadge('video', localVideo.current.id);
      } else {
        showMutedBadge('video', localVideo.current.id);
      }
      onMediaToggle('video');
    }
  }

  const disabled = !name || !hasVideo;

  const title = `${room?.hostName}\ninvites you`;

  const onJoin = () => {
    navigate('/call/' + sid, {
      state: {
        name,
        callState: constraintsState,
        audioEnabled,
        videoEnabled,
        e2ee: room.e2ee,
        title: room.title,
      }
    });
  };

  return (
    <>
      <Header
        title={room?.title}
        centered
      >
        <Flex
          h={40}
          alignItems={'center'}
        >
          <span className={styles.smallText}>at the room:</span>&nbsp;<ParticipantsBadge count={room?.count}/>
        </Flex>
      </Header>

      <Container>
        <Flex
          width={'100%'}
          className={styles.container}
        >
          <div className={styles.videoContainer}>
            <div ref={videoContainer}/>

            <div className={styles.videoControls}>
              <VideoControls
                devices={devices}
                videoEnabled={videoEnabled}
                audioEnabled={audioEnabled}
                onDeviceChange={onDeviceSelect}
                toggleAudio={toggleAudio}
                toggleVideo={toggleVideo}
                selectedVideoId={selectedVideoId}
                selectedAudioId={selectedAudioId}
              />
            </div>
          </div>


          <Flex
            className={styles.joinContainer}
          >
            <h1 className={styles.title}>{title}</h1>

            <Input
              label={'Enter your name:'}
              value={name}
              onChange={setName}
              placeholder={'John'}
              icon={FaceIcon}
            />

            <Box mt={'32px'}>
              <p className={styles.label}>Join as a Participant:</p>
            </Box>
            <div
              className={styles.buttonContainer}
            >
              <Button
                onClick={onJoin}
                text={'Free'}
                disabled={disabled}
              />
            </div>

            <Flex
              mt={isMobile ? 24 : 'auto'}
              color="#555555"
              flexDirection="column"
              gap={8}
            >
              {room.e2ee && (
                <p className={styles.text}>
                  {'End-to-end encryption is enabled,\n' +
                    'we recommend using the Google Chrome browser.'}
                </p>
              )}
            </Flex>
          </Flex>
        </Flex>
      </Container>

      <Footer/>
    </>
  );
};

export default observer(JoinParticipant);
