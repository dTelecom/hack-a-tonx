import React from 'react'
import {ParticipantsIcon} from '../../assets'
import styles from './ParticipantsBadge.module.scss'
import classNames from 'classnames'

const ParticipantsBadge = ({count, isCall}) => {
  return (
    <div className={classNames(styles.container, isCall && styles.isCall)}>
      <img
        src={ParticipantsIcon}
        alt={'person icon'}
      />
      {count}
    </div>
  )
}

export default ParticipantsBadge