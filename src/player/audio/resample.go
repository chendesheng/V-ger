package audio

import (
	. "player/libav"
)

func channelLayout(cnt int) int64 {
	if cnt == 1 {
		return AV_CH_LAYOUT_MONO
	} else {
		return AV_CH_LAYOUT_STEREO
	}
}

func resampleFrame(resampleCtx AVAudioResampleContext, frame AVFrame, codecCtx *AVCodecContext) AVObject {
	channelLayout := channelLayout(codecCtx.Channels())

	// println("resample in:", int64(frame.ChannelLayout()), int64(frame.Format()), int64(frame.SampleRate()))
	// println("resample out:", channelLayout, AV_SAMPLE_FMT_S16, int64(codecCtx.SampleRate()))

	resampleCtxObj := resampleCtx.Object()
	resampleCtxObj.SetOptInt("in_channel_layout", int64(frame.ChannelLayout()), 0)
	resampleCtxObj.SetOptInt("in_sample_fmt", int64(frame.Format()), 0)
	resampleCtxObj.SetOptInt("in_sample_rate", int64(frame.SampleRate()), 0)
	resampleCtxObj.SetOptInt("out_channel_layout", channelLayout, 0)
	resampleCtxObj.SetOptInt("out_sample_fmt", AV_SAMPLE_FMT_S16, 0)
	resampleCtxObj.SetOptInt("out_sample_rate", int64(codecCtx.SampleRate()), 0)

	outChannels := GetChannelLayoutNbChannels(uint64(channelLayout))

	if resampleCtx.Open() < 0 {
		println("error initializing libavresample")
		return AVObject{}
	}
	defer resampleCtx.Close()

	osize := GetBytesPerSample(AV_SAMPLE_FMT_S16)
	outSize, outLinesize := AVSampleGetBufferSize(outChannels, frame.NbSamples(), AV_SAMPLE_FMT_S16)
	// println("frame data size:", outSize)

	tmpOut := AVObject{}
	tmpOut.Malloc(outSize)
	outSamples := resampleCtx.Convert(tmpOut, outLinesize, frame.NbSamples(),
		frame.Data(), frame.Linesize(0), frame.NbSamples())
	if outSamples < 0 {
		println("avresample_convert() failed")
		return AVObject{}
	}
	// println("channels:", codecCtx.Channels(), "outchannels:", outChannels)
	tmpOut.SetSize(outSamples * osize * outChannels)
	return tmpOut
}
