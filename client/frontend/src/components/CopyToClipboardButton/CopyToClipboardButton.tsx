import styles from "../../pages/Call/Call.module.scss"
import React, {useEffect, useRef, useState} from "react"
import {ChainGreenIcon, ChainIcon, GreenTickIcon, WhiteTickIcon} from '../../assets'
import CopyToClipboard from 'react-copy-to-clipboard'
import {useBreakpoints} from "../../hooks/useBreakpoints"

interface IProps {
  text: string
}

export const CopyToClipboardButton = ({text}: IProps) => {
  const timer = useRef<NodeJS.Timeout>()
  const [copied, setCopied] = useState(false)
  const {isMobile} = useBreakpoints()

  useEffect(() => {
    return () => {
      clearTimeout(timer.current)
    }
  }, [])

  function onCopy() {
    if (timer.current) {
      clearTimeout(timer.current)
    }
    setCopied(true)
    timer.current = setTimeout(() => setCopied(false), 2000)
  }

  return (
    <CopyToClipboard
      onCopy={onCopy}
      text={text}
    >
      {isMobile ? (
        <button
          disabled={!text}
          className={styles.inviteButtonMobile}
        >
          <img
            src={copied ? GreenTickIcon : ChainGreenIcon}
            alt={'copy icon'}
          />
        </button>
      ) : (
        <button
          disabled={!text}
          className={styles.inviteButton}
        >
          <img
            src={copied ? WhiteTickIcon : ChainIcon}
            alt={'copy icon'}
          />
          {copied ? 'Copied!' : 'Copy invite link'}
        </button>
      )}
    </CopyToClipboard>
  )
}
