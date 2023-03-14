import React, {useEffect, useState} from 'react';
import axios from 'axios';
import {Header} from '../../components/Header/Header';
import {useNavigate, useParams} from 'react-router-dom';
import ParticipantsBadge from '../../components/ParticipantsBadge/ParticipantsBadge';
import {Flex} from '@chakra-ui/react';
import styles from './JoinModeSelect.module.scss';
import Footer from '../../components/Footer/Footer';
import {JoinCard} from '../../components/JoinCard/JoinCard';
import {JoinScreenshotParticipant, JoinScreenshotViewer} from '../../assets';
import {Container} from '../../components/Container/Container';
import {useBreakpoints} from '../../hooks/useBreakpoints';

export const JoinModeSelect = () => {
  const {isMobile} = useBreakpoints();
  const navigate = useNavigate();
  const {sid} = useParams();
  const [room, setRoom] = useState();

  const loadRoom = async () => {
    axios.post('https://app.dmeet.org/api/room/info', {sid})
      .then((response) => {
        setRoom(response.data);
      })
      .catch(e => {
        console.error(e);
        navigate('/');
      });
  };

  useEffect(() => {
    void loadRoom();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    if (room) {
      if (room.viewerPrice === '') {
        navigate(`/join/participant/${sid}`, {state: {room}});
      }
      if (room.participantPrice === '') {
        navigate(`/join/viewer/${sid}`, {state: {room}});
      }
    }
  }, [navigate, room, sid]);

  if (!room) return null;

  return (
    <>
      <Header
        title={room?.title}
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
          mt={isMobile ? 40 : 80}
          pb={isMobile ? 40 : 24}
          flexDirection={'column'}
          alignItems={'center'}
        >
          <h1 className={styles.title}>
            {room.hostName}{isMobile ? <br/> : ' '}invites you
          </h1>

          <h3 className={styles.subtitle}>You can join as a viewer or as a participant:</h3>

          <Flex
            mt={isMobile ? 24 : 40}
            mb={40}
            gap={isMobile ? 12 : 30}
            flexDirection={isMobile ? 'column' : 'row'}
          >
            <JoinCard
              img={JoinScreenshotViewer}
              text={'View only'}
              buttonText={'Viewer'}
              onClick={() => navigate(`/join/viewer/${sid}`, {state: {room}})}
            />

            <JoinCard
              img={JoinScreenshotParticipant}
              text={'Speak and view'}
              buttonText={'Participant'}
              onClick={() => navigate(`/join/participant/${sid}`, {state: {room}})}
            />
          </Flex>
          {room.e2ee && (
            <p className={styles.smallGreyText}>
              End-to-end encryption is enabled,<br/>
              we recommend using the Google Chrome browser.
            </p>
          )}
        </Flex>
      </Container>

      <Footer/>
    </>
  );
};