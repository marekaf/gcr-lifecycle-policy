package utils

// TranslateVolumeType translates volumeType like "io1" to Pricing API value like "Provisioned IOPS"
func TranslateVolumeType(volumeType string) string {

	// https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ebs-volume-types.html

	typesMap := map[string]string{
		"io1":      "Provisioned IOPS",
		"gp2":      "General Purpose",
		"sc1":      "Cold HDD",
		"st1":      "Throughput Optimized HDD",
		"standard": "Magnetic",
	}

	return typesMap[volumeType]
}
