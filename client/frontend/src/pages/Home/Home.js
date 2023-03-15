import React, {useCallback, useEffect, useMemo, useRef, useState} from 'react';
import {Header} from '../../components/Header/Header';
import styles from './Home.module.scss';
import {observer} from 'mobx-react';
import {useNavigate} from 'react-router-dom';
import {Container} from '../../components/Container/Container';
import {Box, Flex} from '@chakra-ui/react';
import Input from '../../components/Input/Input';
import VideoControls from '../../components/VideoControls/VideoControls';
import Footer from '../../components/Footer/Footer';
import {createVideoElement, hideMutedBadge, showMutedBadge} from '../Call/utils';
import {useMediaConstraints} from '../../hooks/useMediaConstraints';
import {CustomCheckbox} from '../../components/Checkbox/CustomCheckbox';
import {loadDevices} from '../../utils/loadDevices';
import {FaceIcon, KeyIcon} from '../../assets';
import {useBreakpoints} from '../../hooks/useBreakpoints';
import {Button} from '../../components/Button/Button';

const Home = () => {
  const {isMobile} = useBreakpoints();
  const navigate = useNavigate();
  const [hasVideo, setHasVideo] = useState(false);
  const [devices, setDevices] = useState([]);
  const [values, setValues] = useState({
    viewer: true,
    viewerPrice: '0',
    e2ee: true,
    participant: true,
    participantPrice: '0',
    roomName: '',
    name: '',
  });

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
  const videoContainer = useRef();
  const localVideo = useRef();

  useEffect(() => {
    void loadMedia(constraints);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const onChange = useCallback((key, value) => {
    setValues(prev => ({...prev, [key]: value}));
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

  const disabled = useMemo(() => {
    // one of two modes should be selected
    const participantOrViewer = values.participant || values.viewer;
    return !values.name || !values.roomName || !hasVideo || !participantOrViewer;
  }, [values, hasVideo]);

  const title = 'Create a Web3 Video Room';

  const getCallState = () => {
    return {
      callState: constraintsState,
      audioEnabled,
      videoEnabled,
      ...values,
      title: values.roomName,
    };
  };

  const onCreateMeeting = async () => {
    navigate('/call', {
      state: getCallState()
    });
  };

  return (
    <>
      <Header
        title={isMobile ? undefined : title}
        centered={isMobile}
      />

      {isMobile && (
        <h1 className={styles.title}>{title}</h1>
      )}

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

          <Flex maxWidth={isMobile ? 'initial' : '420px'}>
            <Flex
              className={styles.joinContainer}
            >
              <Input
                value={values.name}
                onChange={(value) => onChange('name', value)}
                placeholder={'John'}
                icon={FaceIcon}
                label={'Enter your name:'}
              />

              <Box my={'24px'}>
                <Input
                  value={values.roomName}
                  onChange={(value) => onChange('roomName', value)}
                  placeholder={'What is Web3?'}
                  icon={KeyIcon}
                  label={'Enter a room name:'}
                />
              </Box>

              <Box
                width={'100%'}
                marginTop="24px"
              >
                <CustomCheckbox
                  label={'End-to-end encryption (E2EE)'}
                  checked={values.e2ee}
                  setChecked={(checked) => onChange('e2ee', checked)}
                />
              </Box>

              <Box
                marginTop={isMobile ? '24px' : '56px'}
                width={'100%'}
                boxSizing={'border-box'}
                className={styles.buttonContainer}
              >
                <Button
                  onClick={onCreateMeeting}
                  text={'Create a Video Room'}
                  disabled={disabled}
                />
              </Box>

            </Flex>
          </Flex>
        </Flex>
      </Container>

      <Footer/>
    </>
  );
};

export default observer(Home);
