package audio

import (
	. "player/libav"
)

func channelLayout(cnt int) int64 {
	if cnt == 1 {
		return int64(GetChannelLayout("mono"))
	} else {
		return int64(GetChannelLayout("stereo"))
	}
}

func resampleFrame(resampleCtx AVAudioResampleContext, frame AVFrame, codecCtx *AVCodecContext) AVObject {
	channelLayout := channelLayout(codecCtx.Channels())

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
	outSize, outLinesize := AVSampleGetBufferSize(outChannels, frame.NbSamples(), frame.Format())
	// println("frame data size:", outSize)

	tmpOut := AVObject{}
	tmpOut.Malloc(outSize)
	outSamples := resampleCtx.Convert(tmpOut, outLinesize, frame.NbSamples(),
		frame.Data(), frame.Linesize(0), frame.NbSamples())
	if outSamples < 0 {
		println("avresample_convert() failed")
		return AVObject{}
	}
	// println("channels:", codecCtx.Channels())
	tmpOut.SetSize(outSamples * osize * outChannels)
	return tmpOut
}
