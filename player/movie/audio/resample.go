package audio

import (
	"log"
	. "vger/player/libav"
)

func channelLayout(cnt int) int64 {
	if cnt == 1 {
		return AV_CH_LAYOUT_MONO
	} else {
		return AV_CH_LAYOUT_STEREO
	}
}

func resampleFrame(resampleCtx AVAudioResampleContext, frame AVFrame, codecCtx *AVCodecContext) AVObject {
	inChannelLayout := int64(frame.ChannelLayout())
	if inChannelLayout == 0 {
		inChannelLayout = channelLayout(codecCtx.Channels())
	}

	outChannelLayout := int64(AV_CH_LAYOUT_STEREO)
	// log.Print("resample in:", int64(inChannelLayout), int64(frame.Format()), int64(frame.SampleRate()))
	// log.Print("resample out:", outChannelLayout, AV_SAMPLE_FMT_S16, int64(codecCtx.SampleRate()))

	resampleCtxObj := resampleCtx.Object()
	resampleCtxObj.SetOptInt("in_channel_layout", int64(inChannelLayout), 0)
	resampleCtxObj.SetOptInt("in_sample_fmt", int64(frame.Format()), 0)
	resampleCtxObj.SetOptInt("in_sample_rate", int64(frame.SampleRate()), 0)
	resampleCtxObj.SetOptInt("out_channel_layout", outChannelLayout, 0)
	resampleCtxObj.SetOptInt("out_sample_fmt", AV_SAMPLE_FMT_S16, 0)
	resampleCtxObj.SetOptInt("out_sample_rate", int64(codecCtx.SampleRate()), 0)

	outChannels := 2 //GetChannelLayoutNbChannels(uint64(outChannelLayout))

	if resampleCtx.Open() < 0 {
		log.Print("error initializing libavresample")
		return AVObject{}
	}
	defer resampleCtx.Close()

	osize := GetBytesPerSample(AV_SAMPLE_FMT_S16)
	outSize, outLinesize := AVSampleGetBufferSize(outChannels, frame.NbSamples(), AV_SAMPLE_FMT_S16)
	// log.Print("frame data size:", outSize)

	tmpOut := AVObject{}
	tmpOut.Malloc(outSize)
	outSamples := resampleCtx.Convert(tmpOut, outLinesize, frame.NbSamples(),
		frame.Data(), frame.Linesize(0), frame.NbSamples())
	if outSamples < 0 {
		log.Print("avresample_convert() failed")
		return AVObject{}
	}
	// log.Print("channels:", codecCtx.Channels(), "outchannels:", outChannels)
	tmpOut.SetSize(outSamples * osize * outChannels)
	return tmpOut
}
