package audio

import (
	. "player/libav"
)

func resampleFrame(resampleCtx AVAudioResampleContext, frame AVFrame, channels int) AVObject {
	resampleCtxObj := resampleCtx.Object()
	resampleCtxObj.SetOptInt("in_channel_layout", int64(frame.ChannelLayout()), 0)
	resampleCtxObj.SetOptInt("in_sample_fmt", int64(frame.Format()), 0)
	resampleCtxObj.SetOptInt("in_sample_rate", int64(frame.SampleRate()), 0)
	resampleCtxObj.SetOptInt("out_channel_layout", int64(GetChannelLayout("stereo")), 0)
	resampleCtxObj.SetOptInt("out_sample_fmt", AV_SAMPLE_FMT_S16, 0)
	resampleCtxObj.SetOptInt("out_sample_rate", int64(frame.SampleRate()), 0)

	if resampleCtx.Open() < 0 {
		println("error initializing libavresample")
		return AVObject{}
	}
	defer resampleCtx.Close()

	osize := GetBytesPerSample(AV_SAMPLE_FMT_S16)
	outSize, outLinesize := AVSampleGetBufferSize(channels, frame.NbSamples(), frame.Format())
	// println("frame data size:", outSize)

	tmpOut := AVObject{}
	tmpOut.Malloc(outSize)
	// tmpOut := make([]byte, outSize)
	outSamples := resampleCtx.Convert(tmpOut, outLinesize, frame.NbSamples(),
		frame.Data(), frame.Linesize(0), frame.NbSamples())
	if outSamples < 0 {
		println("avresample_convert() failed")
		return AVObject{}
	}
	tmpOut.SetSize(outSamples * osize * 2)
	// defer tmpOut.Free()
	// return tmpOut.Bytes()
	return tmpOut
}
