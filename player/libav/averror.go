package libav

const (
	AVERROR_BSF_NOT_FOUND      = (-0x39acbd08) ///< Bitstream filter not found
	AVERROR_DECODER_NOT_FOUND  = (-0x3cbabb08) ///< Decoder not found
	AVERROR_DEMUXER_NOT_FOUND  = (-0x32babb08) ///< Demuxer not found
	AVERROR_ENCODER_NOT_FOUND  = (-0x3cb1ba08) ///< Encoder not found
	AVERROR_EOF                = (-0x5fb9b0bb) ///< End of file
	AVERROR_EXIT               = (-0x2bb6a7bb) ///< Immediate exit was requested; the called function should not be restarted
	AVERROR_FILTER_NOT_FOUND   = (-0x33b6b908) ///< Filter not found
	AVERROR_INVALIDDATA        = (-0x3ebbb1b7) ///< Invalid data found when processing input
	AVERROR_MUXER_NOT_FOUND    = (-0x27aab208) ///< Muxer not found
	AVERROR_OPTION_NOT_FOUND   = (-0x2bafb008) ///< Option not found
	AVERROR_PATCHWELCOME       = (-0x3aa8beb0) ///< Not yet implemented in Libav, patches welcome
	AVERROR_PROTOCOL_NOT_FOUND = (-0x30adaf08) ///< Protocol not found
	AVERROR_STREAM_NOT_FOUND   = (-0x2dabac08) ///< Stream not found
	AVERROR_BUG                = (-0x5fb8aabe) ///< Bug detected, please report the issue
	AVERROR_UNKNOWN            = (-0x31b4b1ab) ///< Unknown error, typically from an external library
	AVERROR_EXPERIMENTAL       = (-0x2bb2afa8) ///< Requested feature is flagged experimental. Set strict_std_compliance if you really want to use it.
)
